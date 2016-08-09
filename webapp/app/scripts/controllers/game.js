'use strict';

/**
 * @ngdoc function
 * @name gdApp.controller:GameCtrl
 * @description
 * # GameCtrl
 * Controller of the gdApp
 */
angular.module('gdApp')
  .controller('GameCtrl', ['$scope', '$routeParams', 'game', 'dataService', function ($scope, $routeParams, game, dataService) {

    $scope.debug = true;
    $scope.dart = {Val: 20, Pos: 3};
    $scope.vals = [];
    for (var i = 1; i <= 25; i++) {
      if (i <= 20 || i === 25) {
        if (i === 25) {
          $scope.vals.push({label: 'Bull\'s eye', val: 25});
        } else {
          $scope.vals.push({label: '' + i, val: i});
        }
      }
    }
    $scope.pos = [{label: '', pos: 1}, {label: 'Double', pos: 2}, {label: 'Triple', pos: 3}];
    $scope.alerts = [];

    $scope.closeAlert = function (index) {
      $scope.alerts.splice(index, 1);
    };

    $scope.game = game;

    $scope.getDarts = function (p) {
      if ($scope.game.State.CurrentPlayer === p) {
        return new Array($scope.game.State.CurrentDart + 1);
      }
      return [];
    };

    $scope.addPlayer = function () {
      var name = $scope.newPlayer.name;
      dataService
        .addPlayer($routeParams.id, name)
        .then(
          function (success) {
            if (success) {
             // refresh();
              delete $scope.newPlayer;
            } else {
              $scope.alerts.push({type: 'danger', msg: 'An error occurs'});
            }
          },
          function (failure) {
            $scope.alerts.push({type: 'danger', msg: failure});
          }
        );
    };

    $scope.sendDart = function (sector, pos) {
      dataService
        .sendDart($routeParams.id, sector, pos)
        .then(
          function (state) {
            $scope.game.State = state;
          },
          function (failure) {
            $scope.alerts.push({type: 'danger', msg: failure});
          }
        );
    };

    var cancelled = false;
   
    // join websocket
    dataService
        .joinGame($routeParams.id)
        .then(
            function (ws) {
            console.log("got a ws to bind")
            ws.onmessage = function(event) {
              console.log("got something from space", event)
              $scope.game.State = JSON.parse(event.data);
              $scope.$apply()
            }
          },
          function (failure) {
            $scope.alerts.push({type: 'danger', msg: failure});
          }
        );

  }]);
