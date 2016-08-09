'use strict';

/**
 * @ngdoc service
 * @name gdApp.dataservice
 * @description
 * # dataservice
 * Service in the gdApp.
 */
function DataApi(cacheService, $q, $http) {

  var cache = cacheService.get('data-cache');

  this.styles = function () {
    return cache.getAndSet('styles', function () {
      var q = $q.defer();
      $http
        .get('/api/styles')
        .then(
          function (res) {
            q.resolve(res.data);
          },
          function (rejection) {
            q.reject(rejection);
          });

      return q.promise;
    });
  };


this.games = function () {

    return $http
        .get('/api/games')
        .then(
          function (res) {
            return res.data;
          });
  };

  this.joinGame = function (game) {
    var q = $q.defer();

    var ws = new WebSocket("ws://localhost:8080/api/games/"+ game +"/ws")
    ws.onopen = function() {
       q.resolve (ws);
    }

    return q.promise;
  }

  this.newGame = function (style,board) {

    var q = $q.defer();

    $http
      .post('/api/games', {'Style': style, 'Board': board})
      .then(
        function (response) {
          console.log(response);
          if (response.status === 201) {
            q.resolve(response.data.id);
          } else if (response.data.error) {
            q.reject(response.data.error);
          } else {
            q.reject(response.statusText);
          }
        },
        function (rejection) {
          console.log(rejection);
          if (rejection.data.error) {
            q.reject(rejection.data.error);
          } else if (rejection.statusText) {
            q.reject(rejection.statusText);
          } else {
            q.reject(rejection);
          }
        }
      );

    return q.promise;
  };

  this.game = function (gameId) {
    var q = $q.defer();

    $http
      .get('/api/games/' + gameId)
      .then(
        function (response) {
          q.resolve(response.data.game);
        },
        function (rejection) {
          q.reject(rejection);
        }
      );

    return q.promise;
  };

  this.addPlayer = function (gameId, name) {

    var q = $q.defer();

    $http
      .post('/api/games/' + gameId + '/players', {Name: name})
      .then(
        function (response) {
          console.log(response);
          if (response.status === 201) {
            q.resolve(true);
          } else if (response.data.error) {
            q.reject(response.data.error);
          } else {
            q.reject(response.statusText);
          }
        },
        function (rejection) {
          console.log(rejection);
          if (rejection.data.error) {
            q.reject(rejection.data.error);
          } else if (rejection.statusText) {
            q.reject(rejection.statusText);
          } else {
            q.reject(rejection);
          }
        }
      );

    return q.promise;
  };

  this.sendDart = function (gameId, sector, pos) {
    var q = $q.defer();

    $http
      .post('/api/games/' + gameId + '/darts', {Sector: sector, Multiplier: pos})
      .then(
        function (response) {
          console.log(response);
          if (response.status === 200) {
            q.resolve(response.data.state);
          } else if (response.data.error) {
            q.reject(response.data.error);
          } else {
            q.reject(response.statusText);
          }
        },
        function (rejection) {
          console.log(rejection);
          if (rejection.data.error) {
            q.reject(rejection.data.error);
          } else if (rejection.statusText) {
            q.reject(rejection.statusText);
          } else {
            q.reject(rejection);
          }
        }
      );

    return q.promise;
  };

}

angular.module('gdApp').service('dataService', ['cacheService', '$q', '$http', DataApi]);
