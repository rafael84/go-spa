'use strict';

angular.module('app.storage.media')
    .controller('NewCtrl',
        function($state, FileUploader, Account, Flash, Media, Location, MediaType) {
            var vm = this;
            vm.title = 'media.title.new';
            vm.error = null;
            vm.uploader = new FileUploader({
                headers: Account.getAuthorizationHeader(),
                url: "/api/v1/media/upload"
            });
            vm.media = {};
            vm.save = function(valid) {
                vm.uploader.onSuccessItem = function(fileItem, response, status, headers) {
                    vm.media.path = response;
                    Media.add(vm.media)
                        .then(function success(response) {
                            Flash.show('Media ' + vm.media.name + ' created!');
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
                });
            Media.getMediaTypes()
                .then(function success(response) {
                    vm.mediatypes = response;
                });
        });
