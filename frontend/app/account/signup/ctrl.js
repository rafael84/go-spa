(function() {
    'use strict';
    angular.module('app.account.signup')
        .controller('SignUpCtrl', ['$state', 'Account', 'Flash', SignUpCtrl]);

    function SignUpCtrl($state, Account, Flash) {
        var vm = this;
        vm.user = {};
        vm.error = null;
        Account.getRoles()
            .then(function success(roles) {
                vm.roles = roles;
            });
        vm.register = function register(valid) {
            if (!valid) {
                return;
            }
            vm.user.role = vm.user.role.id;
            Account.signUp(vm.user)
                .then(function success(response) {
                    Flash.show('Thanks for registering!');
                    $state.go('home');
                })
                .catch(function error(response) {
                    vm.error = response.data.error;
                });
        }
    }
})();
