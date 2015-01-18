(function() {
    'use strict';
    angular.module('app.account.signin', [
            'ui.router',
            'app.main',
            'app.account'
        ])
        .config(['$stateProvider', Config]);

    function Config($stateProvider) {
        $stateProvider
            .state('signin', {
                url: '/signin',
                controller: 'SignInCtrl as signin',
                templateUrl: 'app/account/signin/form.html'
            });
    }
})();
