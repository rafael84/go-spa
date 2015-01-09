'use strict';

angular.module("app", [
    'angular-jwt',

    'app.main',
    'app.home',
    'app.account',
    'app.signup',
    'app.signin'
])

.config(function Config($httpProvider, $compileProvider, jwtInterceptorProvider) {
    jwtInterceptorProvider.tokenGetter = function(store) {
        return store.get('token')
    }
    $httpProvider.interceptors.push('jwtInterceptor');
    // $compileProvider.debugInfoEnabled(false);
})

.run(function Run($rootScope, $state, $stateParams) {
    $rootScope.$state = $state;
    $rootScope.$stateParams = $stateParams;
});
