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
        .controller('LocationNewCtrl', ['$state', 'Flash', 'Location', LocationNewCtrl])
        .controller('LocationEditCtrl', ['$state', 'Flash', 'Location', 'location', LocationEditCtrl]);

    function Config($stateProvider) {
        $stateProvider
            .state('location', {
                abstract: true,
                url: '/location',
                template: '<ui-view/>',
                resolve: {
                    Location: 'Location'
                }
            })
            .state('location.list', {
                url: '/list',
                templateUrl: 'app/location/location.list.tmpl.html',
                controller: 'LocationListCtrl as vm',
                resolve: {
                    locations: function(Location) {
                        return Location.getAll();
                    }
                }
            })
            .state('location.new', {
                url: '/new',
                templateUrl: 'app/location/location.new.tmpl.html',
                controller: 'LocationNewCtrl as vm'
            })
            .state('location.edit', {
                url: '/edit/:locationId',
                templateUrl: 'app/location/location.edit.tmpl.html',
                controller: 'LocationEditCtrl as vm',
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
        return {
            getAll: getAll,
            getById: getById,
            remove: remove,
            add: add,
            edit: edit
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

    function LocationNewCtrl($state, Flash, Location) {
        var vm = this;
        vm.error = null;
        vm.location = {};
        vm.save = function(valid) {
            Location.add(vm.location)
                .then(function success(response) {
                    Flash.show('Location ' + vm.location.name + ' created!');
                    $state.go('location.list');
                })
                .catch(function error(response) {
                    vm.error = response;
                });
        }
    }

    function LocationEditCtrl($state, Flash, Location, location) {
        var vm = this;
        vm.error = null;
        vm.location = location;
        vm.save = function(valid) {
            Location.edit(vm.location)
                .then(function success(response) {
                    Flash.show('Location ' + vm.location.name + ' updated!');
                    $state.go('location.list');
                })
                .catch(function error(response) {
                    vm.error = response;
                });
        }
    }
})();
