'use strict';

angular.module('app.account.profile', [
    'ngSanitize',
    'ui.router',
    'ui.select',
    'angular-jwt',
    'angular-storage',
    'app.main'
]).config(function($stateProvider) {
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
});
