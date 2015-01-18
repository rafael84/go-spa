(function() {
    'use strict';
    angular.module('app.storage.media', [
            'ui.router',
            'ngDialog',
            'angular-storage',
            'angular-jwt',
            'app.account',
            'app.main',
            'app.storage.mediatype',
            'app.storage.location'
        ])
        .config(['$stateProvider', Config])

    function Config($stateProvider) {
        $stateProvider
            .state('media', {
                abstract: true,
                url: '/media',
                template: '<ui-view/>',
                resolve: {
                    Media: 'Media'
                }
            })
            .state('media.list', {
                url: '/list',
                templateUrl: 'app/storage/media/list.html',
                controller: 'MediaListCtrl as vm',
                resolve: {
                    medias: function(Media) {
                        return Media.getAll();
                    }
                }
            })
            .state('media.new', {
                url: '/new',
                templateUrl: 'app/storage/media/form.html',
                controller: 'MediaNewCtrl as vm'
            })
            .state('media.edit', {
                url: '/edit/:id',
                templateUrl: 'app/storage/media/form.html',
                controller: 'MediaEditCtrl as vm',
                resolve: {
                    media: function($stateParams, Media) {
                        return Media.getById($stateParams.id);
                    }
                }
            });
    }
})();
