'use strict';

angular.module('app.account.group', [
    'ui.router',
    'ngDialog',
    'angular-storage',
    'angular-jwt',
    'app.main'
]).config(function($stateProvider) {
    $stateProvider
        .state('group', {
            abstract: true,
            url: '/group',
            template: '<ui-view/>',
            resolve: {
                Group: 'Group'
            }
        })
        .state('group.list', {
            url: '/list',
            templateUrl: 'app/account/group/list.html',
            controller: 'GroupListCtrl as vm',
            resolve: {
                groups: function(Group) {
                    return Group.getAll();
                }
            }
        })
        .state('group.new', {
            url: '/new',
            templateUrl: 'app/account/group/form.html',
            controller: 'GroupNewCtrl as vm'
        })
        .state('group.edit', {
            url: '/edit/:groupId',
            templateUrl: 'app/account/group/form.html',
            controller: 'GroupEditCtrl as vm',
            resolve: {
                group: function($stateParams, Group) {
                    return Group.getById($stateParams.groupId);
                }
            }
        });
});
