'use strict';

describe('Filter: sector', function () {

  // load the filter's module
  beforeEach(module('gdApp'));

  // initialize a new instance of the filter before each test
  var sector;
  beforeEach(inject(function ($filter) {
    sector = $filter('sector');
  }));

  it('should return the input prefixed with "sector filter:"', function () {
    var text = 'angularjs';
    expect(sector(text)).toBe('sector filter: ' + text);
  });

});
