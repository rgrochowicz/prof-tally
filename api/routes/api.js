var express = require('express');
var Promise = require('bluebird');
var pg = require('pg');
var _ = require('lodash');

var connectionObject = {
	user: process.env.POSTGRES_USER,
	database: process.env.POSTGRES_DATABASE,
	password: process.env.POSTGRES_PASSWORD,
	host: process.env.POSTGRES_HOST,
	post: process.env.POSTGRES_PORT
};

var router = express.Router();
/* Get list of classes */
router.get('/classes', function(req, res) {
	pg.connect(connectionObject, function(err, client, done) {
		if(err) {
			res.json({error: err});
			return;
		}

		//this should probably be cached, as it's not likely to change
		client.query('SELECT distinct(title), subject || course_num AS course_num, subject FROM courses ORDER BY title', function(err, result) {
			done();

			if(err) {
				res.json({error: err});
				return;
			}
			res.json(result.rows);
		});
	});
});

/* Get schedules based on the titles */
router.get('/schedule', function(req, res) {
	pg.connect(connectionObject, function(err, client, done) {
		if(err) {
			res.json({error: err});
			return;
		}
		var classDetails = req.query.classTitles.map(function(title) {

			//get titles, crns, and class times
			return Promise.promisify(client.query, client)("SELECT courses.title, courses.crn, course_times.weekday, to_char(course_times.start_time, 'HH24:MI') AS start, to_char(course_times.length, 'HH24:MI') AS length FROM course_times, courses WHERE course_times.course_crn = courses.crn AND courses.title = $1;", [title]);
		});
		Promise.all(classDetails).then(function(details) {

			//process for presentation
			//so much processing is going to be done here
			//TODO: do actual processing
			var vals = _(details).pluck('rows').flatten().groupBy('weekday').value();

			res.json(vals);
			done();
		})
	});
	console.log(req.query);
});

module.exports = router;
