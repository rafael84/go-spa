(function() {
    'use strict';
    angular.module('app.media', [
            'ui.router',
            'ngDialog',
            'angular-storage',
            'angular-jwt',
            'app.main',
            'app.mediaType',
            'app.location'
        ])
        .config(Config)
        .factory('Media', ['$http', '$q', 'Location', 'MediaType', Media])
        .controller('MediaListCtrl', ['ngDialog', 'medias', 'Media', 'Flash', MediaListCtrl])
        .controller('MediaNewCtrl', ['$state', 'Flash', 'Media', 'Location', 'MediaType', MediaNewCtrl])
        .controller('MediaEditCtrl', ['$state', 'FileUploader', 'Flash', 'Media', 'Location', 'MediaType', 'media', MediaEditCtrl]);

    function Config($stateProvider) {
        $stateProvider
            .state('media', {
                abstract: true,
                url: '/media',
                template: '<ui-view/>',
                resolve: {
                    Media: 'Media'
                }
            })
            .state('media.list', {
                url: '/list',
                templateUrl: 'app/media/media.list.tmpl.html',
                controller: 'MediaListCtrl as vm',
                resolve: {
                    medias: function(Media) {
                        return Media.getAll();
                    }
                }
            })
            .state('media.new', {
                url: '/new',
                templateUrl: 'app/media/media.new.tmpl.html',
                controller: 'MediaNewCtrl as vm'
            })
            .state('media.edit', {
                url: '/edit/:mediaId',
                templateUrl: 'app/media/media.edit.tmpl.html',
                controller: 'MediaEditCtrl as vm',
                resolve: {
                    media: function($stateParams, Media) {
                        return Media.getById($stateParams.mediaId);
                    }
                }
            });
    }

    function Media($http, $q, Location, MediaType) {
        function getAll() {
            var deferred = $q.defer();
            $http.get("/api/v1/media")
                .then(function success(response) {
                    deferred.resolve(response.data);
                })
                .catch(function error(response) {
                    deferred.reject(response.data.error);
                });
            return deferred.promise;
        }

        function getById(id) {
            var deferred = $q.defer();
            $http.get("/api/v1/media/" + id)
                .then(function success(response) {
                    deferred.resolve(response.data);
                })
                .catch(function error(response) {
                    deferred.reject(response.data.error);
                });
            return deferred.promise;
        }

        function remove(media) {
            var deferred = $q.defer();
            $http.delete("/api/v1/media/" + media.id)
                .then(function success(response) {
                    deferred.resolve(response.data);
                })
                .catch(function error(response) {
                    deferred.reject(response.data.error);
                });
            return deferred.promise;
        }

        function add(media) {
            var deferred = $q.defer();

            media.locationId = media.location.id;
            media.mediaTypeId = media.mediaType.id;

            $http.post("/api/v1/media", media)
                .then(function success(response) {
                    deferred.resolve(response.data);
                })
                .catch(function error(response) {
                    deferred.reject(response.data.error);
                });
            return deferred.promise;
        }

        function edit(media) {
            var deferred = $q.defer();

            media.locationId = media.location.id;
            media.mediaTypeId = media.mediaType.id;

            $http.put("/api/v1/media/" + media.id, media)
                .then(function success(response) {
                    deferred.resolve(response.data);
                })
                .catch(function error(response) {
                    deferred.reject(response.data.error);
                });
            return deferred.promise;
        }

        function getLocations() {
            return Location.getAll();
        }

        function getMediaTypes() {
            return MediaType.getAll();
        }
        return {
            getAll: getAll,
            getById: getById,
            remove: remove,
            add: add,
            edit: edit,
            getLocations: getLocations,
            getMediaTypes: getMediaTypes
        }
    }

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

    function MediaNewCtrl($state, Flash, Media, Location, MediaType) {
        var vm = this;
        vm.error = null;
        vm.media = {};
        vm.save = function(valid) {
            Media.add(vm.media)
                .then(function success(response) {
                    Flash.show('Media ' + vm.media.name + ' created!');
                    $state.go('media.list');
                })
                .catch(function error(response) {
                    vm.error = response;
                });
        }
        Media.getLocations()
            .then(function success(response) {
                vm.locations = response;
            });

        Media.getMediaTypes()
            .then(function success(response) {
                vm.mediaTypes = response;
            });
    }

    function MediaEditCtrl($state, FileUploader, Flash, Media, Location, MediaType, media) {
        var vm = this;
        vm.error = null;
        vm.uploader = new FileUploader({
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
                vm.mediaTypes = response;
                vm.media.mediaType = MediaType.findLocal(vm.media.mediaTypeId, vm.mediaTypes);
            });
    }
})();
