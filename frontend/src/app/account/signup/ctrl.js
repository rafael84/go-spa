'use strict';

angular.module('app.account.signup')
    .controller('SignUpCtrl', function($state, Account, Flash) {
        var vm = this;
        vm.user = {};
        vm.error = null;
        vm.register = function register(valid) {
            if (!valid) {
                return;
            }
            Account.signUp(vm.user)
                .then(function success(response) {
                    Flash.show('Thanks for registering!');
                    $state.go('home');
                })
                .catch(function error(response) {
                    vm.error = response.data.error;
                });
        }
    });
