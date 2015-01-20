'use strict';

angular.module('app.account.group')
    .controller('ListCtrl', function(ngDialog, groups, Group, Flash) {
        var vm = this;
        vm.groups = groups;
        vm.deleteDlg = function(group) {
            vm.group = group;
            ngDialog.open({
                template: 'deleteDlgTmpl',
                data: vm
            });
        }
        vm.delete = function() {
            Group.remove(vm.group)
                .then(function success(response) {
                    Group.getAll()
                        .then(function success(response) {
                            vm.groups = response;
                        });
                    Flash.show("Deleted");
                    vm.group = null;
                })
                .catch(function error(response) {
                    Flash.show("Error!");
                });
        }
    });
