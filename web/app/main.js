angular.module('proftallyApp', [
	"localytics.directives",
	"ngResource",
	"proftallyApp.components",
	"ui.router"
])

.config(['$stateProvider', '$urlRouterProvider', function($stateProvider, $urlRouterProvider) {
	$urlRouterProvider.otherwise('/');

	$stateProvider
	.state('index', {
		url: '/',
		templateUrl: '/partials/index.html',
		controller: 'IndexController'
	})
	.state('params', {
		url: '/params',
		templateUrl: '/partials/params.html',
		controller: 'ParamsController'
	})
	.state('make', {
		url: '/make/:classes/:attrs',
		templateUrl: '/partials/make.html',
		controller: 'MakeController',
		resolve: {
			crns: ['$http', function($http) {
				return $http.get('/api/crns');
			}],
			schedules: ['$stateParams', '$http', function($stateParams, $http) {
				return $http({
					method: 'GET',
					url: '/api/schedule',
					params: {
						"classTitles[]": JSON.parse(atob($stateParams.classes)),
						"attrs[]": JSON.parse(atob($stateParams.attrs))
					}
				})
			}]
		}
	})
}])

.controller('IndexController', ['$scope', function($scope) {

}])

.controller('ParamsController', ['$scope', '$resource', '$state', function($scope, $resource, $state) {
	$scope.classes = $resource('/api/classes').query();
	$scope.attrs = $resource('/api/attrs').query();

	$scope.submit = function() {
		$state.go('make', {
			classes: btoa(JSON.stringify(_.pluck($scope.selectedClasses, 'title'))),
			attrs: btoa(JSON.stringify(_.pluck($scope.selectedAttrs, 'short')))
		});
	}
}])

.controller('MakeController', ['$scope', '$http', '$stateParams', 'crns', 'schedules', function($scope, $http, $stateParams, crns, schedules) {

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
	};

	$scope.currentScheduleIndex = 0;
	$scope.schedules = schedules.data;
	$scope.crns = crns.data;

	//returns object with days based off of events
	$scope.transform = function(schedule) {
		if(!schedule) return [];

		var courses = schedule.map(function(e) { return $scope.crns[e]; });

		var events = [];
		courses.forEach(function(course) {
			course.times.forEach(function(time) {
				events.push(_(course).pick(["title", "crn"]).merge(time).value());
			});
		});

		return _(events).groupBy("weekday").value();
	}

	//filters that will cause a schedule to be rejected
	var filters = [
		{weekday: 'F'}, //exclude classes on friday
		//{weekday: 'T', start: 600}, //exclude classes on tuesday that start before 10:00am
		//{weekday: 'R', end: 1200}, //exclude classes on thursday that end after 8:00pm
		{start: 660}, //exclude classes that start before 11:00am
	];



	var filter = function(schedules) {
		return _.reject(schedules, function(sch) {
			return _.any(sch, function(crn) {
				return _.any($scope.crns[crn].times, function(time) {
					return _.any(filters, function(filter) {
						var results = [];
						if(filter.weekday) results.push(time.weekday == filter.weekday);
						if(filter.start) results.push(time.startminutes < filter.start);
						if(filter.end) results.push(time.endminutes > filter.end);
						return _.all(results);
					});
				});
			});
		});
	}

	$scope.schedules = filter($scope.schedules);


	$scope.prevClick = function() {
		$scope.currentScheduleIndex = Math.max(0, $scope.currentScheduleIndex-1);
	}
	$scope.nextClick = function() {
		$scope.currentScheduleIndex = Math.min($scope.schedules.length-1, $scope.currentScheduleIndex+1);
	}

	$scope.$watch("currentScheduleIndex", function() {
		$scope.events = $scope.transform($scope.schedules[$scope.currentScheduleIndex]);
	})
}]);