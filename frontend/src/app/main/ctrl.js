'use strict';

angular.module('app.main')
    .controller('MainCtrl', function($scope, $translate, ngDialog, Account, Flash) {
        var main = this;
        main.getUser = Account.getUser;
        main.isUserSignedIn = Account.isUserSignedIn;
        main.signOut = Account.signOut;
        main.getFlashMessage = Flash.get;
        main.hasFlashMessage = Flash.hasMessage;
        if (!Account.isTokenExpired()) {
            Account.startTokenRenewal();
        }
        main.switchLang = function(newLang) {
            $translate.use(newLang);
            Flash.show(newLang);
        }
    });
