'use strict';

angular.module('app.account.group')
    .controller('EditCtrl', function($state, Flash, Group, group) {
        var vm = this;
        vm.error = null;
        vm.group = group;
        vm.save = function(valid) {
            Group.edit(vm.group)
                .then(function success(response) {
                    Flash.show('Group ' + vm.group.name + ' updated!');
                    $state.go('group.list');
                })
                .catch(function error(response) {
                    vm.error = response;
                });
        }
    });
