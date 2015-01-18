(function() {
    'use strict';
    angular.module('app.account.signin')
        .controller('SignInCtrl', function SignInCtrl($state, Account, Flash) {
            var signin = this;
            signin.user = {};
            signin.error = null;
            signin.authenticate = function authenticate() {
                Account.signIn(signin.user)
                    .then(function success(response) {
                        Flash.show('Welcome back!');
                        $state.go('home');
                    })
                    .catch(function error(response) {
                        signin.error = response.data.error;
                    });
            }
        })
})();
