API
===

Provides endpoints for schedule creation and stats.

Endpoints:
-------------
####/schedule
    Returns a list of possible schedules based off of course titles.

Post data:
```
courseTitles[]=Course+1&courseTitles[]=Course+2
```
Response:
```json
[{
	"courses": [
		"title": "A Course Title",
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
    Returns the available courses with subjects, titles, and numbers.

Response:
```json
[{
	"subject": "CS",
	"title": "Class Title",
	"course_num": "CS01000"
},...]
```