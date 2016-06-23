'use strict';

/**
 * @ngdoc filter
 * @name gdApp.filter:sector
 * @function
 * @description
 * # sector
 * Filter in the gdApp.
 */
angular.module('gdApp')
  .filter('sector', function () {
    return function (input) {
      if (input.Val === 0) {
        return 'Out of space';
      }
      var prefix = '';
      if(input.Pos === 3) {
        prefix = 'Triple ';
      } else if (input.Pos === 2) {
        prefix = 'Double ';
      }
      if (input.Val === 25) {
        return prefix + 'Bull\'s Eye';
      }
      return prefix + input.Val;
    };
  });
