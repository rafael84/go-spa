(function() {
    'use strict';
    angular.module('app', [
            'ngSanitize',
            'angular-jwt',
            'ui.select',
            'angularFileUpload',
            //
            'app.main',
            'app.home',
            'app.account',
            'app.account.profile',
            'app.account.resetpassword',
            'app.account.signup',
            'app.account.signin',
            'app.account.group',
            'app.storage.location',
            'app.storage.mediatype',
            'app.storage.media'
        ])
        .config(['$httpProvider', '$compileProvider', 'jwtInterceptorProvider', 'uiSelectConfig', Config])
        .run(['$rootScope', '$state', '$stateParams', 'Flash', Run]);

    function Config($httpProvider, $compileProvider, jwtInterceptorProvider, uiSelectConfig) {
        jwtInterceptorProvider.tokenGetter = function(store) {
            return store.get('token')
        }
        $httpProvider.interceptors.push('jwtInterceptor');
        // $compileProvider.debugInfoEnabled(false);
        uiSelectConfig.theme = 'bootstrap';
    }

    function Run($rootScope, $state, $stateParams, Flash) {
        $rootScope.$state = $state;
        $rootScope.$stateParams = $stateParams;
        $rootScope.$on('$stateChangeError',
            function(event, toState, toParams, fromState, fromParams, error) {
                event.preventDefault();
                Flash.show('Unable to access the requested location.');
                $state.go('home');
            });
    }
})();
