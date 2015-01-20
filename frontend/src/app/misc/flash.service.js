'use strict';

angular.module('app.misc')
    .service('Flash', function($translate, $timeout) {
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
    });
