'use strict';

describe('Controller: GamectrlCtrl', function () {

  // load the controller's module
  beforeEach(module('gdApp'));

  var GamectrlCtrl,
    scope;

  // Initialize the controller and a mock scope
  beforeEach(inject(function ($controller, $rootScope) {
    scope = $rootScope.$new();
    GamectrlCtrl = $controller('GamectrlCtrl', {
      $scope: scope
      // place here mocked dependencies
    });
  }));

  it('should attach a list of awesomeThings to the scope', function () {
    expect(GamectrlCtrl.awesomeThings.length).toBe(3);
  });
});
