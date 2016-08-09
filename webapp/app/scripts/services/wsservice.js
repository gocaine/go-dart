'use strict';

function WebSocketService($q, $log) {

  this.open = function open(gameId) {
    var ws = new WebSocket("http://localhost:8080/games/" + gameId + "/ws")
    ws.onopen = function(event) {
      $log.info("yo !")
    }
  };
}

angular.module('gdApp').service('wsService', ['$q', '$log', WebSocketService]);
