'use strict';

angular.module('app.storage.location', [
    'ui.router',
    'ngDialog',
    'angular-storage',
    'angular-jwt',
    'app.main'
]).config(function($stateProvider) {
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
            templateUrl: 'app/storage/location/list.html',
            controller: 'ListCtrl as vm',
            resolve: {
                locations: function(Location) {
                    return Location.getAll();
                }
            }
        })
        .state('location.new', {
            url: '/new',
            templateUrl: 'app/storage/location/form.html',
            controller: 'NewCtrl as vm'
        })
        .state('location.edit', {
            url: '/edit/:locationId',
            templateUrl: 'app/storage/location/form.html',
            controller: 'EditCtrl as vm',
            resolve: {
                location: function($stateParams, Location) {
                    return Location.getById($stateParams.locationId);
                }
            }
        });
});
