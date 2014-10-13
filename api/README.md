API
===

Provides endpoints for schedule creation and stats.

Endpoints:
-------------
####/schedule
Post data:
```
classNames[]=Class1&classNames[]=Class2
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