'use strict';

angular.module('app.account', [
    'ui.router',
    'angular-jwt',
    'angular-storage'
])

.config(function Config($stateProvider) {
    $stateProvider
        .state('resetPassword', {
            url: '/resetPassword',
            controller: 'ResetPasswordCtrl as vm',
            templateUrl: 'app/account/reset-password.tmpl.html'
        });
})

.factory('Account', function Account($http, $q, $interval, $rootScope, store, jwtHelper) {

    var account = this;

    account.renewalPromise = null;

    account.startTokenRenewal = function() {
        account.renewalPromise = $interval(account.renewToken, 30000); // TODO: get this interval in millis from server
    };

    account.stopTokenRenewal = function() {
        if (angular.isDefined(account.renewalPromise)) {
            $interval.cancel(account.renewalPromise);
            account.renewalPromise = undefined;
        }
    };

    account.renewToken = function() {
        return $http.post(
            "/api/v1/accounts/token/renew"
        ).then(function success(response) {
            store.set('token', response.data.token);
            $rootScope.$broadcast('tokenRenewed', response); // TODO: is it really necessary?
        }).catch(function error(response) {
            account.stopTokenRenewal();
            $rootScope.$broadcast('tokenNotRenewed', response) // TODO: is it really necessary?
        });
    };

    account.isTokenExpired = function() {
        var token = store.get('token');
        return token == null || jwtHelper.isTokenExpired(token);
    };

    account.signIn = function(user) {
        return $http.post(
            "/api/v1/accounts/user/signin", user
        ).then(function success(response) {
            store.set('token', response.data.token);
            account.startTokenRenewal();
            $rootScope.$broadcast('user:signedIn', response);
        });
    };

    account.signOut = function() {
        store.remove('token');
        account.stopTokenRenewal();
        $rootScope.$broadcast('user:signedOut', null); // TODO: is it really necessary?
    };

    account.signUp = function(user) {
        var deferred = $q.defer();
        if (user.password !== user.passwordAgain) {
            deferred.reject({
                data: {
                    error: 'Passwords does not match'
                }
            });
        } else {
            return $http.post('/api/v1/accounts/user/signup', user);
        }
        return deferred.promise;
    };

    account.isUserSignedIn = function() {
        return !account.isTokenExpired();
    };

    account.getUser = function() {
        if (account.isTokenExpired()) {
            return null;
        }
        var token = store.get('token');
        var claims = jwtHelper.decodeToken(token);
        return claims.user;
    };

    account.resetPassword = function(user) {
        return $http.post('/api/v1/accounts/user/resetPassword', user);
    };

    return {
        signIn: account.signIn,
        isUserSignedIn: account.isUserSignedIn,
        isTokenExpired: account.isTokenExpired,
        startTokenRenewal: account.startTokenRenewal,
        signOut: account.signOut,
        signUp: account.signUp,
        getUser: account.getUser,
        resetPassword: account.resetPassword
    };
})

.controller('ResetPasswordCtrl', function ResetPasswordCtrl($state, Account) {
    var vm = this;

    vm.error = null;

    vm.send = function send(user) {
        Account.resetPassword(user)
            .then(function success(response) {
                $state.go('home');
            })
            .catch(function error(response) {
                vm.error = response.data.error;
            });
    }
})
