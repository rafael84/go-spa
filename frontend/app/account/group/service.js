(function() {
    'use strict';
    angular.module('app.account.group')
        .factory('Group', ['$http', '$q', Group]);

    function Group($http, $q) {
        function getAll() {
            var deferred = $q.defer();
            $http.get("/api/v1/group")
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
            $http.get("/api/v1/group/" + id)
                .then(function success(response) {
                    deferred.resolve(response.data);
                })
                .catch(function error(response) {
                    deferred.reject(response.data.error);
                });
            return deferred.promise;
        }

        function remove(group) {
            var deferred = $q.defer();
            $http.delete("/api/v1/group/" + group.id)
                .then(function success(response) {
                    deferred.resolve(response.data);
                })
                .catch(function error(response) {
                    deferred.reject(response.data.error);
                });
            return deferred.promise;
        }

        function add(group) {
            var deferred = $q.defer();
            $http.post("/api/v1/group", group)
                .then(function success(response) {
                    deferred.resolve(response.data);
                })
                .catch(function error(response) {
                    deferred.reject(response.data.error);
                });
            return deferred.promise;
        }

        function edit(group) {
            var deferred = $q.defer();
            $http.put("/api/v1/group/" + group.id, group)
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
})();
