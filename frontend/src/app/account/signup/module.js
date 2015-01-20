'use strict';

angular.module('app.account.signup', [
    'ui.router',
    'app.main',
    'app.account'
]).config(function($stateProvider) {
    $stateProvider.state('signup', {
        url: '/signup',
        controller: 'SignUpCtrl as vm',
        templateUrl: 'app/account/signup/form.html'
    });
});
