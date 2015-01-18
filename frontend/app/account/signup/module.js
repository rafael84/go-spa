(function() {
    'use strict';
    angular.module('app.account.signup', [
            'ui.router',
            'app.main',
            'app.account'
        ])
        .config(['$stateProvider', Config]);

    function Config($stateProvider) {
        $stateProvider.state('signup', {
            url: '/signup',
            controller: 'SignUpCtrl as vm',
            templateUrl: 'app/account/signup/form.html'
        });
    }
})();
