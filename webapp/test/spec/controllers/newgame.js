'use strict';

describe('Controller: NewgamectrlCtrl', function () {

  // load the controller's module
  beforeEach(module('gdApp'));

  var NewgamectrlCtrl,
    scope;

  // Initialize the controller and a mock scope
  beforeEach(inject(function ($controller, $rootScope) {
    scope = $rootScope.$new();
    NewgamectrlCtrl = $controller('NewgamectrlCtrl', {
      $scope: scope
      // place here mocked dependencies
    });
  }));

  it('should attach a list of awesomeThings to the scope', function () {
    expect(NewgamectrlCtrl.awesomeThings.length).toBe(3);
  });
});
