(function() {
    'use strict';
    angular.module('app.storage.media')
        .factory('Media', ['$http', '$q', 'Location', 'MediaType', Media]);

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
            media.mediatypeId = media.mediatype.id;
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
            media.mediatypeId = media.mediatype.id;
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
})();
