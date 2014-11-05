package main

import (
	_ "database/sql"
	"encoding/json"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
	"time"
)

var (
	db     *sqlx.DB
	c      redis.Conn
	crnMap map[string]Course
)

type CourseEntry struct {
	Title   string
	Crn     string
	Weekday string
	Start   string
	Length  string
}

type Course struct {
	Crn   string
	Title string
	Times []CourseTime
}

type CourseTime struct {
	Weekday string
	Start   string
	Length  string
}

type CrnSchedule struct {
	Crns []string
}

type TrueSchedule struct {
	Courses []Course
}

//Populate the translation map from crn to course
func populateTranslations() {
	courses := []CourseEntry{}
	db.Select(&courses, `select
	courses.title,
	courses.crn,
	course_times.weekday,
	to_char(course_times.start_time, 'HH24:MI') as start,
	to_char(course_times.length, 'HH24:MI') as length
from
	course_times,
	courses
where
	course_times.course_crn = courses.crn`)

	crnMap = map[string]Course{}
	for _, course := range courses {
		if _, ok := crnMap[course.Crn]; ok {
			existing := crnMap[course.Crn]
			existing.Times = append(crnMap[course.Crn].Times, CourseTime{
				Weekday: course.Weekday,
				Start:   course.Start,
				Length:  course.Length,
			})
		} else {
			crnMap[course.Crn] = Course{
				Title: course.Title,
				Crn:   course.Crn,
				Times: []CourseTime{},
			}
		}
	}
}

//Translate an array of array of crns to schedules
func translateSchedules(schedules []CrnSchedule) []TrueSchedule {
	result := []TrueSchedule{}
	for _, schedule := range schedules {
		newSchedule := TrueSchedule{}
		for _, crn := range schedule.Crns {
			newSchedule.Courses = append(newSchedule.Courses, crnMap[crn])
		}
		result = append(result, newSchedule)
	}

	return result
}

func getNext(crns []string, titles []string) []CrnSchedule {

	//when there are not titles left, return the crns collected
	if len(titles) == 0 {
		return []CrnSchedule{
			CrnSchedule{
				Crns: crns,
			},
		}
	}

	nextTitle := titles[0]

	crnArgs := []string{}
	for _, crn := range crns {
		crnArgs = append(crnArgs, fmt.Sprintf("crn:%s", crn))
	}

	//get all available courses based on set intersections
	entries, err := redis.Strings(c.Do("SINTER", redis.Args{}.AddFlat(crnArgs).Add(fmt.Sprintf("title:%s", nextTitle))...))
	if err != nil {
		log.Fatalln(err)
	}


	finalResult := []CrnSchedule{}
	for _, entry := range entries {
		results := getNext(append(crns, entry), titles[1:])
		for _, result := range results {
			finalResult = append(finalResult, result)
		}
	}

	return finalResult

}

func main() {

	var err error

	//connect to postgres
	db, err = sqlx.Connect("postgres", "user=tally dbname=tally password=tally sslmode=disable")
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()

	//connect to redis
	c, err = redis.Dial("tcp", ":6379")
	if err != nil {
		log.Fatalln(err)
	}
	defer c.Close()

	titles := []string{"FOUNDATIONS COMP SCI", "OBJ-ORIENT PRGRM/DATA ABSTR", "LINEAR ALGEBRA", "COLLEGE COMPOSITION II", "PUBLIC SPEAKING"}

	startTime := time.Now()

	populateTranslations()
	crnSchedules := getNext([]string{}, titles)
	trueSchedules := translateSchedules(crnSchedules)

	jsonBytes, err := json.Marshal(trueSchedules)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println(time.Since(startTime))

	log.Println(len(crnSchedules))
	log.Println(len(jsonBytes))

}
