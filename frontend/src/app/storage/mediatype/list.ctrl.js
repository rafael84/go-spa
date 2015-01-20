'use strict';

angular.module('app.storage.mediatype')
    .controller('MediatypeListCtrl', function(ngDialog, mediatypes, MediaType, Flash) {
        var vm = this;
        vm.mediatypes = mediatypes;
        vm.deleteDlg = function(mediatype) {
            vm.mediatype = mediatype;
            ngDialog.open({
                template: 'deleteDlgTmpl',
                data: vm
            });
        }
        vm.delete = function() {
            MediaType.remove(vm.mediatype)
                .then(function success(response) {
                    MediaType.getAll()
                        .then(function success(response) {
                            vm.mediatypes = response;
                        });
                    Flash.show("Deleted");
                    vm.mediatype = null;
                })
                .catch(function error(response) {
                    Flash.show("Error!");
                });
        }
    });
