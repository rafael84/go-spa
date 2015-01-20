'use strict';

angular.module('app.account.profile')
    .controller('ProfileCtrl', function(user, Profile, roles, Flash) {
        var vm = this;
        vm.user = user;
        vm.roles = roles;
        vm.update = function(valid) {
            vm.user.role = vm.user.role.id;
            Profile.update(vm.user)
                .then(function success(response) {
                    Flash.show('Your profile has been updated.');
                })
                .catch(function error(response) {
                    vm.error = response;
                });
        }
    });
