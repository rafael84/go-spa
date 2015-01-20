'use strict';

angular.module('app.i18n', ['pascalprecht.translate'])
    .config(function($translateProvider) {
        $translateProvider.useLoader('translationLoader', {});
        $translateProvider.preferredLanguage('en');
    })
    .factory('translationLoader', function($http, $q) {
        return function(options) {
            var deferred = $q.defer();
            $http.get('assets/translations/' + options.key + '.json')
                .then(function success(response) {
                    deferred.resolve(response.data);
                });
            return deferred.promise;
        }
    });
