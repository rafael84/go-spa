(function() {
    'use strict';
    angular.module('app.group', [
            'ui.router',
            'ngDialog',
            'angular-storage',
            'angular-jwt',
            'app.main'
        ])
        .config(Config)
        .factory('Group', ['$http', '$q', Group])
        .controller('GroupListCtrl', ['ngDialog', 'groups', 'Group', 'Flash', GroupListCtrl])
        .controller('GroupDetailCtrl', ['group', GroupDetailCtrl]);

    function Config($stateProvider) {
        $stateProvider
            .state('group', {
                url: '/groups',
                templateUrl: 'app/group/group.list.tmpl.html',
                controller: 'GroupListCtrl as vm',
                resolve: {
                    Group: 'Group',
                    groups: function(Group) {
                        return Group.getAll();
                    }
                }
            })
            .state('group.detail', {
                url: '/:groupId',
                templateUrl: 'app/group/group.detail.tmpl.html',
                controller: 'GroupDetailCtrl as vm',
                resolve: {
                    group: function($stateParams, Group) {
                        return Group.getById($stateParams.groupId);
                    }
                }
            });
    }

    function Group($http, $q) {
        function getAll() {
            var deferred = $q.defer();
            $http.get("/api/v1/account/group")
                .then(function success(response) {
                    deferred.resolve(response.data);
                })
                .catch(function error(response) {
                    deferred.reject(response.data.error);
                });
            return deferred.promise;
        }

        function getById(id) {
            var deferred = $q.defer();
            $http.get("/api/v1/account/group/" + id)
                .then(function success(response) {
                    deferred.resolve(response.data);
                })
                .catch(function error(response) {
                    deferred.reject(response.data.error);
                });
            return deferred.promise;
        }

        function remove(group) {
            var deferred = $q.defer();
            $http.delete("/api/v1/account/group/" + group.id)
                .then(function success(response) {
                    deferred.resolve(response.data);
                })
                .catch(function error(response) {
                    deferred.reject(response.data.error);
                });
            return deferred.promise;
        }
        return {
            getAll: getAll,
            getById: getById,
            remove: remove
        }
    }

    function GroupListCtrl(ngDialog, groups, Group, Flash) {
        var vm = this;
        vm.groups = groups;
        vm.deleteDlg = function(group) {
            vm.group = group;
            ngDialog.open({
                template: 'deleteDlgTmpl',
                data: vm
            });
        }
        vm.delete = function() {
            Group.remove(vm.group)
                .then(function success(response) {
                    Group.getAll()
                        .then(function success(response) {
                            vm.groups = response;
                        });
                    Flash.show("Deleted");
                    vm.group = null;
                })
                .catch(function error(response) {
                    Flash.show("Error!");
                });
        }
    }

    function GroupDetailCtrl(group) {
        var vm = this;
        vm.group = group;
    }
})();
