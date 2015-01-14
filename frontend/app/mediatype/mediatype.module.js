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
        .controller('MediaTypeNewCtrl', ['$state', 'Flash', 'MediaType', MediaTypeNewCtrl])
        .controller('MediaTypeEditCtrl', ['$state', 'Flash', 'MediaType', 'mediaType', MediaTypeEditCtrl]);

    function Config($stateProvider) {
        $stateProvider
            .state('mediaType', {
                abstract: true,
                url: '/mediatype',
                template: '<ui-view/>',
                resolve: {
                    MediaType: 'MediaType'
                }
            })
            .state('mediaType.list', {
                url: '/list',
                templateUrl: 'app/mediatype/mediatype.list.tmpl.html',
                controller: 'MediaTypeListCtrl as vm',
                resolve: {
                    mediaTypes: function(MediaType) {
                        return MediaType.getAll();
                    }
                }
            })
            .state('mediaType.new', {
                url: '/new',
                templateUrl: 'app/mediatype/mediatype.new.tmpl.html',
                controller: 'MediaTypeNewCtrl as vm'
            })
            .state('mediaType.edit', {
                url: '/edit/:mediaTypeId',
                templateUrl: 'app/mediatype/mediatype.edit.tmpl.html',
                controller: 'MediaTypeEditCtrl as vm',
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
            $http.get("/api/v1/mediatype")
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
            $http.get("/api/v1/mediatype/" + id)
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
            $http.delete("/api/v1/mediatype/" + mediaType.id)
                .then(function success(response) {
                    deferred.resolve(response.data);
                })
                .catch(function error(response) {
                    deferred.reject(response.data.error);
                });
            return deferred.promise;
        }

        function add(mediaType) {
            var deferred = $q.defer();
            $http.post("/api/v1/mediatype", mediaType)
                .then(function success(response) {
                    deferred.resolve(response.data);
                })
                .catch(function error(response) {
                    deferred.reject(response.data.error);
                });
            return deferred.promise;
        }

        function edit(mediaType) {
            var deferred = $q.defer();
            $http.put("/api/v1/mediatype/" + mediaType.id, mediaType)
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
            remove: remove,
            add: add,
            edit: edit
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

    function MediaTypeNewCtrl($state, Flash, MediaType) {
        var vm = this;
        vm.error = null;
        vm.mediaType = {};
        vm.save = function(valid) {
            MediaType.add(vm.mediaType)
                .then(function success(response) {
                    Flash.show('MediaType ' + vm.mediaType.name + ' created!');
                    $state.go('mediaType.list');
                })
                .catch(function error(response) {
                    vm.error = response;
                });
        }
    }

    function MediaTypeEditCtrl($state, Flash, MediaType, mediaType) {
        var vm = this;
        vm.error = null;
        vm.mediaType = mediaType;
        vm.save = function(valid) {
            MediaType.edit(vm.mediaType)
                .then(function success(response) {
                    Flash.show('MediaType ' + vm.mediaType.name + ' updated!');
                    $state.go('mediaType.list');
                })
                .catch(function error(response) {
                    vm.error = response;
                });
        }
    }
})();
