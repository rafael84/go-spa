'use strict';

angular.module('app.storage.mediatype')
    .service('MediaType', function($http, $q) {
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

        function remove(mediatype) {
            var deferred = $q.defer();
            $http.delete("/api/v1/mediatype/" + mediatype.id)
                .then(function success(response) {
                    deferred.resolve(response.data);
                })
                .catch(function error(response) {
                    deferred.reject(response.data.error);
                });
            return deferred.promise;
        }

        function add(mediatype) {
            var deferred = $q.defer();
            $http.post("/api/v1/mediatype", mediatype)
                .then(function success(response) {
                    deferred.resolve(response.data);
                })
                .catch(function error(response) {
                    deferred.reject(response.data.error);
                });
            return deferred.promise;
        }

        function edit(mediatype) {
            var deferred = $q.defer();
            $http.put("/api/v1/mediatype/" + mediatype.id, mediatype)
                .then(function success(response) {
                    deferred.resolve(response.data);
                })
                .catch(function error(response) {
                    deferred.reject(response.data.error);
                });
            return deferred.promise;
        }

        function findLocal(id, mediatypes) {
            for (var i = 0; i < mediatypes.length; i++) {
                var mediatype = mediatypes[i];
                if (mediatype.id == id) {
                    return mediatype;
                }
            }
            return null;
        }
        return {
            getAll: getAll,
            getById: getById,
            remove: remove,
            add: add,
            edit: edit,
            findLocal: findLocal
        }
    });
