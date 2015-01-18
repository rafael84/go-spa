(function() {
    'use scrict';
    angular.module('app.account.profile')
        .factory('Profile', ['$http', '$q', Profile]);

    function Profile($http, $q) {
        function get() {
            var deferred = $q.defer();
            $http.get('/api/v1/account/user/profile')
                .then(function success(response) {
                    deferred.resolve(response.data);
                })
                .catch(function error(response) {
                    deferred.reject(response.data.error);
                });
            return deferred.promise;
        }

        function update(user) {
            var deferred = $q.defer();
            $http.put('/api/v1/account/user/profile', user)
                .then(function success(response) {
                    deferred.resolve(response.data);
                })
                .catch(function error(response) {
                    deferred.reject(response.data.error);
                });
            return deferred.promise;
        }
        return {
            get: get,
            update: update
        }
    }
})();
