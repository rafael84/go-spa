'use strict';

angular.module('app.account.group')
    .controller('NewCtrl', function($state, Flash, Group) {
        var vm = this;
        vm.error = null;
        vm.group = {};
        vm.save = function(valid) {
            Group.add(vm.group)
                .then(function success(response) {
                    Flash.show('Group ' + vm.group.name + ' created!');
                    $state.go('group.list');
                })
                .catch(function error(response) {
                    vm.error = response;
                });
        }
    });
