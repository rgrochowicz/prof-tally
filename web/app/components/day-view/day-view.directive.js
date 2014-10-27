angular.module("proftallyApp.components", [])

.directive("dayView", ['$parse', function($parse) {
	return {
		restrict: 'E',
		scope: {
			'events': '=?dayViewEvents',
			'title': '=?dayViewTitle',
			'hourBegin': '@dayViewHourBegin',
			'hourEnd': '@dayViewHourEnd'
		},
		templateUrl: 'app/components/day-view/day-view.template.html',
		link: function($scope, elm, attrs) {
			$scope.hourRange = _.range(parseInt($scope.hourBegin || 0), parseInt($scope.hourEnd || 24));

			$scope.getTimeString = function(hour) {
				return moment().hour(hour).format('ha');
			}

			$scope.getEventStyle = function(event) {

				//this is totally scientific and taken from chrome dev tools
				var pixelsPer30Mins = 22;

				var startingTime = moment.duration(event.start),
					length = moment.duration(event.length);

				var timeToTop = function(time) {
					var minuteTimes = time.asMinutes();
					return (minuteTimes / 30) * pixelsPer30Mins;
				}

				var timeOffset = function() {
					//calculates the needed offset for a different hourBegin than 0

					return parseInt($scope.hourBegin) * -2 * pixelsPer30Mins;
				}

				//sets the top of the event and its size
				var eventTop = timeToTop(startingTime) + timeOffset(),
					eventBottom = -eventTop - timeToTop(length);


				//get a unique color based off of the CRN
				var crnHash = murmurhash3_32_gc(event.crn.toString(), 31);
				var color = "hsl("+(crnHash & 0x168) +",70%,40%)";

				return {
					top: eventTop + "px",
					bottom: eventBottom + "px",
					backgroundColor: color
				}
			}

			$scope.getEventTimeString = function(event) {
				var startingTime = moment(event.start, 'HH:mm'),
					length = moment.duration(event.length);

				var endingTime = startingTime.clone().add(length);

				return startingTime.format('h:mma') + " - " + endingTime.format('h:mma');
			}

		}
	}
}]);