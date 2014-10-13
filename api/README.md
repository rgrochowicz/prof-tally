API
===

Provides endpoints for schedule creation and stats.

Endpoints:
-------------
####/schedule
    Returns a list of possible schedules based off of course numbers.

Post data:
```
courseNumbers[]=12315&courseNumbers[]=21342
```
Response:
```json
[{
	"courses": [
		"name": "A Course Name",
		"crn": 1,
		"times": [{
			"day": "M",
			"start": "12:00",
			"length": "1:15"
		}]
    ]
}]
```

####/courses
    Returns the available courses with names and numbers.

Response:
```json
[{
	"course_num": 12345,
	"name": "Class Name"
},...]
```