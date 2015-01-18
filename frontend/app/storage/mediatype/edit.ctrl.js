(function() {
    'use strict';
    angular.module('app.storage.mediatype')
        .controller('MediaTypeEditCtrl', ['$state', 'Flash', 'MediaType', 'mediatype', MediaTypeEditCtrl]);

    function MediaTypeEditCtrl($state, Flash, MediaType, mediatype) {
        var vm = this;
        vm.error = null;
        vm.mediatype = mediatype;
        vm.save = function(valid) {
            MediaType.edit(vm.mediatype)
                .then(function success(response) {
                    Flash.show('MediaType ' + vm.mediatype.name + ' updated!');
                    $state.go('mediatype.list');
                })
                .catch(function error(response) {
                    vm.error = response;
                });
        }
    }
})();
