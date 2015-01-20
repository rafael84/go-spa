'use strict';

angular.module('app', [
    'angular-jwt',
    'ui.select',
    'ui.router',
    'angular-storage',
    'app.misc',
    'app.home',
    'app.account.group',
    'app.account.profile',
    'app.account.resetpassword',
    'app.account.signin',
    'app.account.signup',
    'app.storage.location',
    'app.storage.media',
    'app.storage.mediatype'
]).config(function($httpProvider, $compileProvider, jwtInterceptorProvider, uiSelectConfig) {
    jwtInterceptorProvider.tokenGetter = function(store) {
        return store.get('token');
    };
    $httpProvider.interceptors.push('jwtInterceptor');
    // $compileProvider.debugInfoEnabled(false);
    uiSelectConfig.theme = 'bootstrap';
}).run(function($rootScope, $state, $stateParams, Flash) {
    $rootScope.$state = $state;
    $rootScope.$stateParams = $stateParams;
    $rootScope.$on('$stateChangeError',
        function(event, toState, toParams, fromState, fromParams, error) {
            event.preventDefault();
            Flash.show('Unable to access the requested location.');
            $state.go('home');
        });
});
