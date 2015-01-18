(function() {
    'use strict';
    angular.module('app.account.group', [
            'ui.router',
            'ngDialog',
            'angular-storage',
            'angular-jwt',
            'app.main'
        ])
        .config(['$stateProvider', Config]);

    function Config($stateProvider) {
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
    }

    function Group($http, $q) {
        function getAll() {
            var deferred = $q.defer();
            $http.get("/api/v1/group")
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
            $http.get("/api/v1/group/" + id)
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
            $http.delete("/api/v1/group/" + group.id)
                .then(function success(response) {
                    deferred.resolve(response.data);
                })
                .catch(function error(response) {
                    deferred.reject(response.data.error);
                });
            return deferred.promise;
        }

        function add(group) {
            var deferred = $q.defer();
            $http.post("/api/v1/group", group)
                .then(function success(response) {
                    deferred.resolve(response.data);
                })
                .catch(function error(response) {
                    deferred.reject(response.data.error);
                });
            return deferred.promise;
        }

        function edit(group) {
            var deferred = $q.defer();
            $http.put("/api/v1/group/" + group.id, group)
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
            remove: remove,
            add: add,
            edit: edit
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

    function GroupNewCtrl($state, Flash, Group) {
        var vm = this;
        vm.error = null;
        vm.group = {};
        vm.save = function(valid) {
            Group.add(vm.group)
                .then(function success(response) {
                    Flash.show('Group ' + vm.group.name + ' created!');
                    $state.go('group.list');
                })
                .catch(function error(response) {
                    vm.error = response;
                });
        }
    }

    function GroupEditCtrl($state, Flash, Group, group) {
        var vm = this;
        vm.error = null;
        vm.group = group;
        vm.save = function(valid) {
            Group.edit(vm.group)
                .then(function success(response) {
                    Flash.show('Group ' + vm.group.name + ' updated!');
                    $state.go('group.list');
                })
                .catch(function error(response) {
                    vm.error = response;
                });
        }
    }
})();
