'use strict';

angular.module('app.home', [
    'ui.router',
    'ngDialog',
    'angular-storage',
    'angular-jwt',
    'app.main'
]).config(function($stateProvider) {
    $stateProvider
        .state('home', {
            url: '/',
            templateUrl: 'app/home/page.html'
        });
});
