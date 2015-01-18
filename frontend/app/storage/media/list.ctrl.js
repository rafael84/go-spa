(function() {
    'use strict';
    angular.module('app.storage.media')
        .controller('MediaListCtrl', ['ngDialog', 'medias', 'Media', 'Flash', MediaListCtrl]);

    function MediaListCtrl(ngDialog, medias, Media, Flash) {
        var vm = this;
        vm.medias = medias;
        vm.deleteDlg = function(media) {
            vm.media = media;
            ngDialog.open({
                template: 'deleteDlgTmpl',
                data: vm
            });
        }
        vm.delete = function() {
            Media.remove(vm.media)
                .then(function success(response) {
                    Media.getAll()
                        .then(function success(response) {
                            vm.medias = response;
                        });
                    Flash.show("Deleted");
                    vm.media = null;
                })
                .catch(function error(response) {
                    Flash.show("Error!");
                });
        }
    }
})();
