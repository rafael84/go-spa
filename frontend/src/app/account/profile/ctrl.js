'use strict';

angular.module('app.account.profile')
    .controller('ProfileCtrl', function(user, Profile, Flash) {
        var vm = this;
        vm.user = user;
        vm.update = function(valid) {
            Profile.update(vm.user)
                .then(function success(response) {
                    Flash.show('Your profile has been updated.');
                })
                .catch(function error(response) {
                    vm.error = response;
                });
        }
    });
