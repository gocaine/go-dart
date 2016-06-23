'use strict';

/**
 * @ngdoc filter
 * @name gdApp.filter:histo
 * @function
 * @description
 * # histo
 * Filter in the gdApp.
 */
angular.module('gdApp')
  .filter('histo', [function () {
    return function (input) {
      var out = '';
      if (input && Object.keys(input).length > 0) {
        for (var i = 1; i <= 25; i++) {
          if (i <= 20 || i === 25) {
            if (input[i]) {
              var val = '' + i;
              if (i === 25) {
                val = 'BE';
              }
              var suffix = '' + input[i];
              if (input[i] === 3) {
                suffix = 'Open';
              }
              out += val + ' : ' + suffix + '<br>';
            }
          }
        }
      }
      return out;
    };
  }]);
