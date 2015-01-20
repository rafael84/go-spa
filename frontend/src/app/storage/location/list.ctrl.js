'use strict';

angular.module('app.storage.location')
    .controller('ListCtrl', function(ngDialog, locations, Location, Flash) {
        var vm = this;
        vm.locations = locations;
        vm.deleteDlg = function(location) {
            vm.location = location;
            ngDialog.open({
                template: 'deleteDlgTmpl',
                data: vm
            });
        }
        vm.delete = function() {
            Location.remove(vm.location)
                .then(function success(response) {
                    Location.getAll()
                        .then(function success(response) {
                            vm.locations = response;
                        });
                    Flash.show("Deleted");
                    vm.location = null;
                })
                .catch(function error(response) {
                    Flash.show("Error!");
                });
        }
    });
