'use strict';

describe('Filter: histo', function () {

  // load the filter's module
  beforeEach(module('gdApp'));

  // initialize a new instance of the filter before each test
  var histo;
  beforeEach(inject(function ($filter) {
    histo = $filter('histo');
  }));

  it('should return the input prefixed with "histo filter:"', function () {
    var text = 'angularjs';
    expect(histo(text)).toBe('histo filter: ' + text);
  });

});
