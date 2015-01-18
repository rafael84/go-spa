(function() {
    'use strict';
    angular.module('app.i18n', [
            'pascalprecht.translate'
        ])
        .config(function Config($translateProvider) {
            $translateProvider.useLoader('translationLoader', {});
            $translateProvider.preferredLanguage('en');
        })
        .factory('translationLoader', function translationLoader($http, $q) {
            return function(options) {
                var deferred = $q.defer();
                $http.get('app/i18n/' + options.key + '.json')
                    .then(function success(response) {
                        deferred.resolve(response.data);
                    });
                return deferred.promise;
            }
        })
})();
