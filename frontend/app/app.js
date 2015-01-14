'use strict';

angular.module("app", [
    'ngSanitize',
    'angular-jwt',
    'ui.select',
    'angularFileUpload',

    'app.main',
    'app.home',
    'app.account',
    'app.signup',
    'app.signin',
    'app.group',
    'app.location',
    'app.mediaType',
    'app.media'
])

.config(function Config($httpProvider, $compileProvider, jwtInterceptorProvider, uiSelectConfig) {
    jwtInterceptorProvider.tokenGetter = function(store) {
        return store.get('token')
    }
    $httpProvider.interceptors.push('jwtInterceptor');
    // $compileProvider.debugInfoEnabled(false);

    uiSelectConfig.theme = 'bootstrap';
})

.run(function Run($rootScope, $state, $stateParams) {
    $rootScope.$state = $state;
    $rootScope.$stateParams = $stateParams;
});
