'use strict';

angular.module('app.storage.location')
    .controller('EditCtrl', function($state, Flash, Location, location) {
        var vm = this;
        vm.error = null;
        vm.location = location;
        vm.save = function(valid) {
            Location.edit(vm.location)
                .then(function success(response) {
                    Flash.show('Location ' + vm.location.name + ' updated!');
                    $state.go('location.list');
                })
                .catch(function error(response) {
                    vm.error = response;
                });
        }
    });
