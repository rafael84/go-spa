'use strict';

angular.module('app.storage.media')
    .controller('MediaEditCtrl',
        function($state, FileUploader, Account, Flash, Media, Location, MediaType, media) {
            var vm = this;
            vm.title = 'media.title.edit';
            vm.error = null;
            vm.uploader = new FileUploader({
                headers: Account.getAuthorizationHeader(),
                url: "/api/v1/media/upload"
            });
            vm.media = media;
            vm.save = function(valid) {
                vm.uploader.onSuccessItem = function(fileItem, response, status, headers) {
                    vm.media.path = response;
                    Media.edit(vm.media)
                        .then(function success(response) {
                            Flash.show('Media ' + vm.media.name + ' updated!');
                            $state.go('media.list');
                        })
                        .catch(function error(response) {
                            vm.error = response;
                        });
                }
                var item = vm.uploader.queue[0];
                item.upload();
            }
            Media.getLocations()
                .then(function success(response) {
                    vm.locations = response;
                    vm.media.location = Location.findLocal(vm.media.locationId, vm.locations);
                });
            Media.getMediaTypes()
                .then(function success(response) {
                    vm.mediatypes = response;
                    vm.media.mediatype = MediaType.findLocal(vm.media.mediatypeId, vm.mediatypes);
                });
        });
