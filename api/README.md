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
	"courses": [{
		"title": "A Course Title",
		"crn": 1,
		"times": [{
			"day": "M",
			"start": "12:00",
			"length": "1:15"
		}]
    },...]
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


Set up:
-------
1. Install node ~0.11
2. Run `npm install` to install dependencies
3. Make sure postgres is running and accessible.
4. Run `./bin/www` to start the API server on port 3000

Docker:
-------
A docker container is provided.  Make sure to set the `POSTGRES_HOST`, `POSTGRES_USER`, `POSTGRES_PASSWORD`, and `POSTGRES_DATABASE` to their actual values. The API server will listen on port 3000 by default unless the `PORT` environment variable is modified.