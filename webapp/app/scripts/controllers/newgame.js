'use strict';

/**
 * @ngdoc function
 * @name gdApp.controller:NewgamectrlCtrl
 * @description
 * # NewgamectrlCtrl
 * Controller of the gdApp
 */
angular.module('gdApp')
  .controller('NewGameCtrl', ['$scope', '$location', 'dataService', function ($scope, $location, dataService) {

    $scope.alerts = [];

    $scope.closeAlert = function (index) {
      $scope.alerts.splice(index, 1);
    };

    dataService.styles().then(
      function (data) {
        $scope.styles = data.styles;
      },
      function (rejection) {
        $scope.alerts.push({type: 'danger', msg: rejection});
      });

    $scope.newGame = function (style) {
      console.log('New Game of style : ', style);
      dataService.newGame(style).then(
        function (gameId) {
          console.log('new game created with id', gameId);
          $location.url('game/' + gameId);
        },
        function (reject) {
          $scope.alerts.push({type: 'danger', msg: reject});
        }
      );
    };

  }]);
