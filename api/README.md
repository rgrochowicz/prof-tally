API
===

Provides endpoints for schedule creation and stats.

Endpoints:
-------------
1. /schedule
   Post data: classNames[]
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
			}...]
	    ]
	},...]
   ```