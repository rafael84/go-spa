(function() {
    'use strict';
    angular.module('app.mediaType', [
            'ui.router',
            'ngDialog',
            'angular-storage',
            'angular-jwt',
            'app.main'
        ])
        .config(Config)
        .factory('MediaType', ['$http', '$q', MediaType])
        .controller('MediaTypeListCtrl', ['ngDialog', 'mediaTypes', 'MediaType', 'Flash', MediaTypeListCtrl])
        .controller('MediaTypeDetailCtrl', ['mediaType', MediaTypeDetailCtrl]);

    function Config($stateProvider) {
        $stateProvider
            .state('mediaType', {
                url: '/media-types',
                templateUrl: 'app/mediatype/mediatype.list.tmpl.html',
                controller: 'MediaTypeListCtrl as vm',
                resolve: {
                    MediaType: 'MediaType',
                    mediaTypes: function(MediaType) {
                        return MediaType.getAll();
                    }
                }
            })
            .state('mediaType.detail', {
                url: '/:mediaTypeId',
                templateUrl: 'app/mediatype/mediatype.detail.tmpl.html',
                controller: 'MediaTypeDetailCtrl as vm',
                resolve: {
                    mediaType: function($stateParams, MediaType) {
                        return MediaType.getById($stateParams.mediaTypeId);
                    }
                }
            });
    }

    function MediaType($http, $q) {
        function getAll() {
            var deferred = $q.defer();
            $http.get("/api/v1/storage/mediatype")
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
            $http.get("/api/v1/storage/mediatype/" + id)
                .then(function success(response) {
                    deferred.resolve(response.data);
                })
                .catch(function error(response) {
                    deferred.reject(response.data.error);
                });
            return deferred.promise;
        }

        function remove(mediaType) {
            var deferred = $q.defer();
            $http.delete("/api/v1/storage/mediatype/" + mediaType.id)
                .then(function success(response) {
                    deferred.resolve(response.data);
                })
                .catch(function error(response) {
                    deferred.reject(response.data.error);
                });
            return deferred.promise;
        }
        return {
            getAll: getAll,
            getById: getById,
            remove: remove
        }
    }

    function MediaTypeListCtrl(ngDialog, mediaTypes, MediaType, Flash) {
        var vm = this;
        vm.mediaTypes = mediaTypes;
        vm.deleteDlg = function(mediaType) {
            vm.mediaType = mediaType;
            ngDialog.open({
                template: 'deleteDlgTmpl',
                data: vm
            });
        }
        vm.delete = function() {
            MediaType.remove(vm.mediaType)
                .then(function success(response) {
                    MediaType.getAll()
                        .then(function success(response) {
                            vm.mediaTypes = response;
                        });
                    Flash.show("Deleted");
                    vm.mediaType = null;
                })
                .catch(function error(response) {
                    Flash.show("Error!");
                });
        }
    }

    function MediaTypeDetailCtrl(mediaType) {
        var vm = this;
        vm.mediaType = mediaType;
    }
})();
