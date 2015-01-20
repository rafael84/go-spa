'use strict';

angular.module('app.storage.mediatype', [
    'ui.router',
    'ngDialog',
    'angular-storage',
    'angular-jwt',
    'app.main'
]).config(function($stateProvider) {
    $stateProvider
        .state('mediatype', {
            abstract: true,
            url: '/mediatype',
            template: '<ui-view/>',
            resolve: {
                MediaType: 'MediaType'
            }
        })
        .state('mediatype.list', {
            url: '/list',
            templateUrl: 'app/storage/mediatype/list.html',
            controller: 'MediatypeListCtrl as vm',
            resolve: {
                mediatypes: function(MediaType) {
                    return MediaType.getAll();
                }
            }
        })
        .state('mediatype.new', {
            url: '/new',
            templateUrl: 'app/storage/mediatype/form.html',
            controller: 'MediatypeNewCtrl as vm'
        })
        .state('mediatype.edit', {
            url: '/edit/:id',
            templateUrl: 'app/storage/mediatype/form.html',
            controller: 'MediatypeEditCtrl as vm',
            resolve: {
                mediatype: function($stateParams, MediaType) {
                    return MediaType.getById($stateParams.id);
                }
            }
        });
});
