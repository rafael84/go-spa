'use strict';
angular.module("app.main", [
        'ngDialog',
        'app.account',
        'app.translations'
    ])
    .controller("MainCtrl", function MainCtrl($scope, $translate, ngDialog, Account, Flash) {
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
    })
    .factory("Flash", function Flash($translate, $timeout) {
        var flash = this;
        flash.hideScheduled = false;
        flash.message = null;
        flash.show = function show(message, timeout) {
            flash.message = {
                message: message,
                timeout: timeout || 5000
            };
        }
        flash.hide = function hide() {
            flash.message = null;
            flash.hideScheduled = false;
        }
        flash.get = function get() {
            if (flash.message == null) {
                return null;
            }
            if (!flash.hideScheduled) {
                flash.hideScheduled = true;
                $timeout(flash.hide, flash.message.timeout);
            }
            return flash.message.message;
        }
        flash.hasMessage = function hasMessage() {
            return flash.message != null;
        }
        return {
            show: flash.show,
            hide: flash.hide,
            get: flash.get,
            hasMessage: flash.hasMessage
        }
    })
    .filter('propsFilter', function() {
        return function(items, props) {
            var out = [];
            if (angular.isArray(items)) {
                items.forEach(function(item) {
                    var itemMatches = false;
                    var keys = Object.keys(props);
                    for (var i = 0; i < keys.length; i++) {
                        var prop = keys[i];
                        if (props[prop] == undefined) {
                            continue;
                        }
                        var text = props[prop].toLowerCase();
                        if (item[prop].toString().toLowerCase().indexOf(text) !== -1) {
                            itemMatches = true;
                            break;
                        }
                    }
                    if (itemMatches) {
                        out.push(item);
                    }
                });
            } else {
                // Let the output be the input untouched
                out = items;
            }
            return out;
        }
    });
