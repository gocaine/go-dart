'use strict';

function CacheApi($q, $cacheFactory, _) {

  this.get = function get(cacheName) {
    var cache = $cacheFactory.get(cacheName);
    if (!cache) {
      cache = $cacheFactory(cacheName);
      cache.getAndSet = function (key, f) {
        var q = $q.defer();
        if (this.get(key)) {
          q.resolve(this.get(key));
        } else {
          f().then(function (data) {
            q.resolve(this.put(key, data));
          }.bind(this), function (rejection) {
            q.reject(rejection);
          });
        }
        return q.promise;
      }.bind(cache);
    }
    return cache;
  };

  this.clear = function clear() {
    _.each($cacheFactory.info(), function (truc, key) {
      $cacheFactory.get(key).removeAll();
    });
  };
}

angular.module('gdApp').service('cacheService', ['$q', '$cacheFactory', '_', CacheApi]);
