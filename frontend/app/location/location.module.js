(function() {
    'use strict';
    angular.module('app.location', [
            'ui.router',
            'ngDialog',
            'angular-storage',
            'angular-jwt',
            'app.main'
        ])
        .config(Config)
        .factory('Location', ['$http', '$q', Location])
        .controller('LocationListCtrl', ['ngDialog', 'locations', 'Location', 'Flash', LocationListCtrl])
        .controller('LocationDetailCtrl', ['location', LocationDetailCtrl]);

    function Config($stateProvider) {
        $stateProvider
            .state('location', {
                url: '/locations',
                templateUrl: 'app/location/location.list.tmpl.html',
                controller: 'LocationListCtrl as vm',
                resolve: {
                    Location: 'Location',
                    locations: function(Location) {
                        return Location.getAll();
                    }
                }
            })
            .state('location.detail', {
                url: '/:locationId',
                templateUrl: 'app/location/location.detail.tmpl.html',
                controller: 'LocationDetailCtrl as vm',
                resolve: {
                    location: function($stateParams, Location) {
                        return Location.getById($stateParams.locationId);
                    }
                }
            });
    }

    function Location($http, $q) {
        function getAll() {
            var deferred = $q.defer();
            $http.get("/api/v1/storage/location")
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
            $http.get("/api/v1/storage/location/" + id)
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
            $http.delete("/api/v1/storage/location/" + location.id)
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

    function LocationListCtrl(ngDialog, locations, Location, Flash) {
        var vm = this;
        vm.locations = locations;
        vm.deleteDlg = function(location) {
            vm.location = location;
            ngDialog.open({
                template: 'deleteDlgTmpl',
                data: vm
            });
        }
        vm.delete = function() {
            Location.remove(vm.location)
                .then(function success(response) {
                    Location.getAll()
                        .then(function success(response) {
                            vm.locations = response;
                        });
                    Flash.show("Deleted");
                    vm.location = null;
                })
                .catch(function error(response) {
                    Flash.show("Error!");
                });
        }
    }

    function LocationDetailCtrl(location) {
        var vm = this;
        vm.location = location;
    }
})();
