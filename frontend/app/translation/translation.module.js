'use strict';

angular.module('app.translations', [
    'pascalprecht.translate'
])

.config(function Config($translateProvider) {
    $translateProvider.useLoader('translationLoader', {});
    $translateProvider.preferredLanguage('en');
})

.factory('translationLoader', function translationLoader($http, $q) {
    return function(options) {
        var deferred = $q.defer();
        $http.get('app/translation/' + options.key + '.json')
            .then(function success(response) {
                deferred.resolve(response.data);
            });
        return deferred.promise;
    }
})
