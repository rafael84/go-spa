'use strict';

angular.module('app.account')
    .factory('Account', function($http, $q, $interval, $rootScope, store, jwtHelper) {
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
                "/api/v1/account/token/renew"
            ).then(function success(response) {
                store.set('token', response.data.token);
                $rootScope.$broadcast('tokenRenewed', response); // TODO: is it really necessary?
                return response;
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
                "/api/v1/account/user/signin", user
            ).then(function success(response) {
                store.set('token', response.data.token);
                account.startTokenRenewal();
                $rootScope.$broadcast('user:signedIn', response);
                return response;
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
                return $http.post('/api/v1/account/user/signup', user);
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
            return $http.post('/api/v1/account/reset-password', user);
        };
        account.resetPasswordValidateKey = function(key) {
            return $http.post('/api/v1/account/reset-password/validate-key', {
                key: key
            });
        };
        account.changePassword = function(user) {
            var deferred = $q.defer();
            if (user.password !== user.passwordAgain) {
                deferred.reject({
                    data: {
                        error: 'Passwords does not match'
                    }
                });
            } else {
                return $http.post('/api/v1/account/reset-password/complete', user);
            }
            return deferred.promise;
        };
        account.getAuthorizationHeader = function() {
            if (account.isTokenExpired()) {
                return null;
            }
            var token = store.get('token');
            return {
                'Authorization': 'Bearer ' + token
            };
        };
        account.getRoles = function() {
            var deferred = $q.defer();
            $http.get('/api/v1/account/user/role')
                .then(function success(response) {
                    deferred.resolve(response.data);
                })
                .catch(function error(response) {
                    deferred.reject(response);
                });
            return deferred.promise;
        };
        return account;
    });
