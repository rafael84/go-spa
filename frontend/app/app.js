'use strict';

angular.module("app", [
    'angular-jwt',

    'app.translations',
    'app.account',
    'app.home',
    'app.signup',
    'app.signin'
])

.config(function Config($urlRouterProvider, $httpProvider, $compileProvider, jwtInterceptorProvider) {
    $urlRouterProvider.otherwise('/');

    jwtInterceptorProvider.tokenGetter = function(store) {
        return store.get('token')
    }

    $httpProvider.interceptors.push('jwtInterceptor');

    // $compileProvider.debugInfoEnabled(false);
})

.run(function Run($rootScope, $state, $stateParams) {
    $rootScope.$state = $state;
    $rootScope.$stateParams = $stateParams;
})

.controller("MainCtrl", function MainCtrl($scope, $translate, Account) {
    var main = this;

    main.getUser = Account.getUser;
    main.isUserSignedIn = Account.isUserSignedIn;
    main.signOut = Account.signOut;

    if (!Account.isTokenExpired()) {
        Account.startTokenRenewal();
    }

    main.switchLang = function(newLang) {
        $translate.use(newLang);
    }
});
