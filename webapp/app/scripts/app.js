'use strict';

/**
 * @ngdoc overview
 * @name gdApp
 * @description
 * # gdApp
 *
 * Main module of the application.
 */
angular
  .module('gdApp', [
    'ngRoute',
    'ngSanitize',
    'ng-lodash'
  ])
  .config(function ($routeProvider) {
    $routeProvider
      .when('/', {
        redirectTo: '/game/new'
      })
      .when('/game/new', {
        templateUrl: 'views/newgame.html',
        controller: 'NewGameCtrl'
      })
      .when('/game/:id', {
        templateUrl: 'views/game.html',
        controller: 'GameCtrl',
        resolve: {
          game: ['$route', 'dataService', function ($route, dataService) {
            return dataService.game($route.current.params.id);
          }]
        }
      })
      .otherwise({
        redirectTo: '/'
      });
  });
