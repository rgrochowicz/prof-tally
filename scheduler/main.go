package main

import (
	_ "github.com/lib/pq"
    _ "database/sql"
    "github.com/jmoiron/sqlx"
    "log"
    "strings"
    "fmt"
    "encoding/json"
    "time"
)

var (
	db *sqlx.DB
)

type CourseEntry struct {
	Title string
	Crn string
	Weekday string
	Start string
	Length string
}

type Course struct {
	Crn string
	Title string
	Times []CourseTime
}

type CourseTime struct {
	Weekday string
	Start string
	Length string
}

type Schedule struct {
	Courses []Course
}


//Takes a slice of courses and returns it formatted for PostgreSQL use
//Used for the ANY expression
//ex. [1 2 3] -> {1,2,3}
func coursesToCrnList(courses []Course) string {
	crnArray := []string{}
	for _, course := range courses {
		crnArray = append(crnArray, course.Crn)
	}
	return fmt.Sprintf("{%s}", strings.Join(crnArray, ","))
}

//Converts from a slice of courses to a map that maps from crn to
// the relevant entries
func groupByCrn(courses []CourseEntry) map[string][]CourseEntry {
	result := map[string][]CourseEntry{}
	for _, course := range courses {
		if _, ok := result[course.Crn]; ok {
			result[course.Crn] = append(result[course.Crn], course)
		} else {
			result[course.Crn] = []CourseEntry{
				course,
			}
		}
	}

	return result
}

//Takes a query and params and returns the results in a crn to course list map
func queryCourse(query string, params ...interface{}) map[string]Course {
	result := map[string]Course{}

	entries := []CourseEntry{}
	db.Select(&entries, query, params...)
	mapped := groupByCrn(entries)

	for crn, courseEntries := range mapped {

		//format the result based on the api spec
		times := []CourseTime{}
		for _, entry := range courseEntries {
			times = append(times, CourseTime{
				Weekday: entry.Weekday,
				Start: entry.Start,
				Length: entry.Length,
			})
		}

		result[crn] = Course{
			Crn: crn,
			Title: courseEntries[0].Title,
			Times: times,
		}
	}

	return result

}

func getNext(courses []Course, titles []string) []Schedule {

	//when there are no titles left to search, return the list of courses
	if len(titles) == 0 {
		return []Schedule{
			Schedule{
				Courses: courses,
			},
		}
	}

	nextTitle := titles[0]

	//get the courses with a certain title that doesn't overlap an array of crns
	entries := queryCourse(`with taken_times as (
	select weekday, start_time, length
	from course_times
	where course_crn = ANY($1)
), relevant_crns as (
	select crn
	from courses
	where title = $2
), excluded_crns as (
	select course_crn
	from course_times
	where course_crn in (select crn from relevant_crns) and
	exists (
		select 1
		from taken_times
		where 
			(taken_times.start_time, taken_times.length) overlaps (course_times.start_time, course_times.length) and
			weekday = course_times.weekday
	)
)
select
	courses.title,
	courses.crn,
	course_times.weekday,
	to_char(course_times.start_time, 'HH24:MI') as start,
	to_char(course_times.length, 'HH24:MI') as length
from
	courses,
	course_times
where
	course_times.course_crn = courses.crn and
	crn in (select crn from relevant_crns) and
	crn not in (select course_crn from excluded_crns)`, coursesToCrnList(courses), nextTitle)

	finalResult := []Schedule{}
	for _, entry := range entries {
		results := getNext(append(courses, entry), titles[1:])
		for _, result := range results {
			finalResult = append(finalResult, result)
		}
	}

	return finalResult

}

func main() {

	var err error

	//connect to the postgres db
	db, err = sqlx.Connect("postgres", "user=tally dbname=tally password=tally sslmode=disable")
    if err != nil {
        log.Fatalln(err)
    }

    defer db.Close()
    titles := []string{"FOUNDATIONS COMP SCI","OBJ-ORIENT PRGRM/DATA ABSTR","LINEAR ALGEBRA","COLLEGE COMPOSITION II"}
    
    startTime := time.Now()

    //get all possible schedules
    schedules := getNext([]Course{}, titles)

    log.Println(time.Since(startTime))

    jsonBytes, err := json.Marshal(schedules)
    if err != nil {
    	log.Fatalln(err)
    }
    log.Println(len(jsonBytes))
    log.Println(len(schedules))

}