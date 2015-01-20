'use strict';

angular.module('app.account.resetpassword')
    .controller('ResetPasswordCtrl', function($state, $stateParams, Account, Flash) {
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
    });
