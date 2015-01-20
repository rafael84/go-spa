 'use strict';

 angular.module('app.storage.location')
     .factory('Location', function($http, $q) {
         function getAll() {
             var deferred = $q.defer();
             $http.get("/api/v1/location")
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
             $http.get("/api/v1/location/" + id)
                 .then(function success(response) {
                     deferred.resolve(response.data);
                 })
                 .catch(function error(response) {
                     deferred.reject(response.data.error);
                 });
             return deferred.promise;
         }

         function remove(location) {
             var deferred = $q.defer();
             $http.delete("/api/v1/location/" + location.id)
                 .then(function success(response) {
                     deferred.resolve(response.data);
                 })
                 .catch(function error(response) {
                     deferred.reject(response.data.error);
                 });
             return deferred.promise;
         }

         function add(location) {
             var deferred = $q.defer();
             $http.post("/api/v1/location", location)
                 .then(function success(response) {
                     deferred.resolve(response.data);
                 })
                 .catch(function error(response) {
                     deferred.reject(response.data.error);
                 });
             return deferred.promise;
         }

         function edit(location) {
             var deferred = $q.defer();
             $http.put("/api/v1/location/" + location.id, location)
                 .then(function success(response) {
                     deferred.resolve(response.data);
                 })
                 .catch(function error(response) {
                     deferred.reject(response.data.error);
                 });
             return deferred.promise;
         }

         function findLocal(id, locations) {
             for (var i = 0; i < locations.length; i++) {
                 var location = locations[i];
                 if (location.id == id) {
                     return location;
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
