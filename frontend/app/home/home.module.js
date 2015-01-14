(function() {
    'use strict';
    angular.module('app.home', [
            'ui.router',
            'ngDialog',
            'angular-storage',
            'angular-jwt',
            'app.main'
        ])
        .config(Config)
        .controller('HomeCtrl', ['$state', HomeCtrl])

    function Config($stateProvider) {
        $stateProvider
            .state('home', {
                url: '/',
                templateUrl: 'app/home/home.tmpl.html',
                controller: 'HomeCtrl as vm',
            });
    }

    function HomeCtrl($state) {
        var vm = this;
        vm.error = null;
    }
})();
