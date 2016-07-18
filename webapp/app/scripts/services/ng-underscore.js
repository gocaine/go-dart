'use strict';
var lodash = angular.module('ng-lodash', []);

lodash.provider('_', [function _Provider() {

  this.$get = window._;
}]);

lodash.factory('_', function () {
  return window._;
});
