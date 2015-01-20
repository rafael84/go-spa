'use strict';

angular.module('app.storage.location')
    .controller('NewCtrl', function LocationNewCtrl($state, Flash, Location) {
        var vm = this;
        vm.error = null;
        vm.location = {};
        vm.save = function(valid) {
            Location.add(vm.location)
                .then(function success(response) {
                    Flash.show('Location ' + vm.location.name + ' created!');
                    $state.go('location.list');
                })
                .catch(function error(response) {
                    vm.error = response;
                });
        }
    });
