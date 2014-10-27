angular.module('proftallyApp', [
	"localytics.directives",
	"ngResource",
	"proftallyApp.components"
])

.controller('MainCtrl', ['$scope', '$http', '$resource', function($scope, $http, $resource) {

	$scope.classes = $resource('/api/classes').query();

	$scope.days = ["M", "T", "W", "R", "F", "S"];

	$scope.dayNames = {
		"M": "Monday",
		"T": "Tuesday",
		"W": "Wednesday",
		"R": "Thursday",
		"F": "Friday",
		"S": "Saturday"
	};

	$scope.events = {
		"M": [],
		"T": [],
		"W": [],
		"R": [],
		"F": [],
		"S": []
	}

	$scope.make = function() {
		var classTitles = _.pluck($scope.selectedClasses, 'title');

		$http({
			method: 'GET',
			url: '/api/schedule',
			params: {
				"classTitles[]": classTitles
			}
		}).success(function(data, status) {
			$scope.events = data;
		});
	}

}]);