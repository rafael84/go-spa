'use strict';

angular.module('app.account.resetpassword', [
    'ui.router',
    'ui.select',
    'angular-jwt',
    'angular-storage',
    'app.main'
]).config(function($stateProvider) {
    $stateProvider
        .state('resetPassword', {
            url: '/reset-password',
            controller: 'ResetPasswordCtrl as vm',
            templateUrl: 'app/account/resetpassword/step1.html',
            data: {
                step: 1
            }
        })
        .state('resetPasswordStep2', {
            url: '/reset-password/step2/:key',
            controller: 'ResetPasswordCtrl as vm',
            templateUrl: 'app/account/resetpassword/step2.html',
            data: {
                step: 2
            }
        });
});
