(function() {
    'use strict';
    angular.module('app.account.profile', [
            'ui.router',
            'ui.select',
            'angular-jwt',
            'angular-storage',
            'app.main'
        ])
        .config(['$stateProvider', Config]);

    function Config($stateProvider) {
        $stateProvider
            .state('profile', {
                url: '/profile',
                templateUrl: 'app/account/profile/form.html',
                controller: 'ProfileCtrl as vm',
                resolve: {
                    Profile: 'Profile',
                    user: function(Profile) {
                        return Profile.get();
                    },
                    roles: function(Account) {
                        return Account.getRoles();
                    }
                }
            });
    }
})();
