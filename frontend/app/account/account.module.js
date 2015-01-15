'use strict';
angular.module('app.account', [
        'ui.router',
        'angular-jwt',
        'angular-storage',
        'app.main'
    ])
    .factory('Me', ['$http', '$q',
        function Me($http, $q) {
            function get() {
                var deferred = $q.defer();
                $http.get('/api/v1/account/user/me')
                    .then(function success(response) {
                        deferred.resolve(response.data);
                    })
                    .catch(function error(response) {
                        deferred.reject(response.data.error);
                    });
                return deferred.promise;
            }

            function update(user) {
                var deferred = $q.defer();
                $http.put('/api/v1/account/user/me', user)
                    .then(function success(response) {
                        deferred.resolve(response.data);
                    })
                    .catch(function error(response) {
                        deferred.reject(response.data.error);
                    });
                return deferred.promise;
            }
            return {
                get: get,
                update: update
            }
        }
    ])
    .config(function Config($stateProvider) {
        $stateProvider
            .state('me', {
                url: '/me',
                templateUrl: 'app/account/me.tmpl.html',
                controller: 'MeCtrl as vm',
                resolve: {
                    Me: 'Me',
                    user: function(Me) {
                        return Me.get();
                    }
                }
            })
            .state('resetPassword', {
                url: '/reset-password',
                controller: 'ResetPasswordCtrl as vm',
                templateUrl: 'app/account/reset-password.tmpl.html',
                data: {
                    step: 1
                }
            })
            .state('resetPasswordStep2', {
                url: '/reset-password/step2/:key',
                controller: 'ResetPasswordCtrl as vm',
                templateUrl: 'app/account/reset-password-step2.tmpl.html',
                data: {
                    step: 2
                }
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
        return {
            signIn: account.signIn,
            isUserSignedIn: account.isUserSignedIn,
            isTokenExpired: account.isTokenExpired,
            startTokenRenewal: account.startTokenRenewal,
            signOut: account.signOut,
            signUp: account.signUp,
            getUser: account.getUser,
            resetPassword: account.resetPassword,
            resetPasswordValidateKey: account.resetPasswordValidateKey,
            changePassword: account.changePassword,
            getAuthorizationHeader: account.getAuthorizationHeader
        };
    })
    .controller('MeCtrl', ['user', 'Me', 'Flash',
        function MeCtrl(user, Me, Flash) {
            var vm = this;
            vm.user = user;
            vm.update = function(valid) {
                Me.update(vm.user)
                    .then(function success(response) {
                        Flash.show('Your profile has been updated.');
                    })
                    .catch(function error(response) {
                        vm.error = response;
                    });
            }
        }
    ])
    .controller('ResetPasswordCtrl', function ResetPasswordCtrl($state, $stateParams, Account, Flash) {
        var vm = this;
        vm.error = null;
        vm.validKey = null;
        if ($state.current.data.step == 1) {
            vm.send = function send(user) {
                Account.resetPassword(user)
                    .then(function success(response) {
                        Flash.show('Check your email address.');
                        $state.go('home');
                    })
                    .catch(function error(response) {
                        vm.error = response.data.error;
                    });
            }
        }
        if ($state.current.data.step == 2) {
            Account.resetPasswordValidateKey($stateParams.key)
                .then(function success(response) {
                    vm.validKey = response.data;
                    vm.user = {
                        validKey: vm.validKey
                    }
                })
                .catch(function error(response) {
                    vm.error = response.data.error;
                });
            vm.send = function send(user) {
                Account.changePassword(user)
                    .then(function success(response) {
                        Flash.show('Your account has been updated, you can login now.');
                        $state.go('signin');
                    })
                    .catch(function error(response) {
                        vm.error = response.data.error;
                    });
            }
        }
    })
