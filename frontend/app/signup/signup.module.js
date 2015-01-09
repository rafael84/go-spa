'use strict';

angular.module('app.signup', [
    'ui.router',

    'app.main',
    'app.account'
])

.config(function Config($stateProvider) {
    $stateProvider.state('signup', {
        url: '/signup',
        controller: 'SignUpCtrl as signup',
        templateUrl: 'app/signup/signup.tmpl.html'
    });
})

.controller('SignUpCtrl', function SignUpCtrl($state, Account, Flash) {
    var signup = this;

    signup.user = {};
    signup.error = null;

    signup.register = function register(valid) {
        if (!valid) {
            return;
        }

        Account.signUp(signup.user)
            .then(function success(response) {
                Flash.show('Thanks for registering!');
                $state.go('home');
            })
            .catch(function error(response) {
                signup.error = response.data.error;
            });
    }
})
