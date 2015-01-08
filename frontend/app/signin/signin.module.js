'use strict';

angular.module('app.signin', [
    'ui.router',

    'app.account'
])

.config(function Config($stateProvider) {
    $stateProvider
        .state('signin', {
            url: '/signin',
            controller: 'SignInCtrl as signin',
            templateUrl: 'app/signin/signin.tmpl.html'
        });
})


.controller('SignInCtrl', function SignInCtrl($state, Account) {
    var signin = this;

    signin.user = {};
    signin.error = null;

    signin.authenticate = function authenticate() {
        Account.signIn(signin.user)
            .then(function success(response) {
                $state.go('home');
            })
            .catch(function error(response) {
                signin.error = response.data.error;
            });
    }
})
