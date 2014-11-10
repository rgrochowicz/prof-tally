package main

import (
	"net/http"
	"log"
	"github.com/garyburd/redigo/redis"
	_ "github.com/lib/pq"
    _ "database/sql"
    "github.com/jmoiron/sqlx"
    "encoding/json"
    "fmt"
    "time"
	"github.com/nutrun/lentil"
	"os"
)

type CourseDescription struct {
	Title string `json:"title"`
	CourseNum string `json:"course_num"`
	Subject string `json:"subject"`
}

type RawCourseTime struct {
	Title string `json:"title"`
	Crn string `json:"crn"`
	Weekday string `json:"weekday"`
	Start string `json:"start"`
	StartMinutes int `json:"startminutes"`
	Length string `json:"length"`
	LengthMinutes int `json:"lengthminutes"`
	End string `json:"end"`
	EndMinutes int `json:"endminutes"`
}

type CourseTime struct {
	Weekday string `json:"weekday"`
	Start string `json:"start"`
	StartMinutes int `json:"startminutes"`
	Length string `json:"length"`
	LengthMinutes int `json:"lengthminutes"`
	End string `json:"end"`
	EndMinutes int `json:"endminutes"`
}

type Course struct {
	Title string `json:"title"`
	Crn string `json:"crn"`
	Times []CourseTime `json:"times"`
}

type CourseTimes []RawCourseTime

func (c *CourseTimes) GroupByCrns() map[string]Course {
	result := map[string]Course{}

	for _, courseTime := range *c {
		transformed := CourseTime{
			Weekday: courseTime.Weekday,
			Start: courseTime.Start,
			StartMinutes: courseTime.StartMinutes,
			Length: courseTime.Length,
			LengthMinutes: courseTime.LengthMinutes,
			End: courseTime.End,
			EndMinutes: courseTime.EndMinutes,
		}
		if _, ok := result[courseTime.Crn]; !ok {
			result[courseTime.Crn] = Course{
				Title: courseTime.Title,
				Crn: courseTime.Crn,
				Times: []CourseTime{
					transformed,	
				},
			}
		} else {
			course := result[courseTime.Crn]
			course.Times = append(course.Times, transformed)
			result[courseTime.Crn] = course
		}
	}

	return result
}

func main() {

	log.Println("Listening: api")

	pool := &redis.Pool{
        MaxIdle: 3,
        IdleTimeout: 240 * time.Second,
        Dial: func () (redis.Conn, error) {
            c, err := redis.Dial("tcp", os.ExpandEnv("${REDIS_HOST}:${REDIS_PORT}"))
            if err != nil {
                return nil, err
            }
            return c, err
        },
        TestOnBorrow: func(c redis.Conn, t time.Time) error {
            _, err := c.Do("PING")
            return err
        },
    }

	db, err := sqlx.Connect("postgres", os.ExpandEnv("user=${POSTGRES_USER} dbname=${POSTGRES_DATABASE} password=${POSTGRES_PASSWORD} host=${POSTGRES_HOST} port=${POSTGRES_PORT} sslmode=disable"))
    if err != nil {
        log.Fatalln(err)
    }
    defer db.Close()

	lentil.ReaderSize = 12000000

	http.HandleFunc("/api/classes", func(w http.ResponseWriter, r *http.Request) {

		conn := pool.Get()
		defer conn.Close()

		classes, err := redis.Bytes(conn.Do("GET", "cache:classes"))

		if err != nil {
			descriptions := []CourseDescription{}

			err := db.Select(&descriptions, `select
	distinct(title) as title,
	subject || course_num as coursenum,
	subject as subject
from courses
order by title`)

			if err != nil {
				fmt.Fprintf(w, "%s\n", err)
				return
			}

			classes, err = json.Marshal(descriptions)

			if err != nil {
				fmt.Fprintf(w, "%s\n", err)
				return
			}

			_, err = conn.Do("SET", "cache:classes", string(classes))

			if err != nil {
				fmt.Fprintf(w, "%s\n", err)
				return
			}

			conn.Do("EXPIRE", "cache:classes", "60")

		}

		w.Write(classes)
	})

	http.HandleFunc("/api/crns", func(w http.ResponseWriter, r *http.Request) {

		conn := pool.Get()
		defer conn.Close()

		crns, err := redis.Bytes(conn.Do("GET", "cache:crns"))

		if err != nil {
			times := CourseTimes{}

			err := db.Select(&times, `select
	courses.title,
	courses.crn,
	course_times.weekday,
	to_char(course_times.start_time, 'HH24:MI') as start,
	extract(hour from course_times.start_time) * 60 + extract(minute from course_times.start_time) as startminutes,
	to_char(course_times.length, 'HH24:MI') as length,
	extract(hour from course_times.length) * 60 + extract(minute from course_times.length) as lengthminutes,
	to_char(course_times.start_time + course_times.length, 'HH24:MI') as end,
	extract(hour from course_times.start_time) * 60 + extract(minute from course_times.start_time) + extract(hour from course_times.length) * 60 + extract(minute from course_times.length) as endminutes
from
	course_times,
	courses
where
	course_times.course_crn = courses.crn and
	course_times.weekday is not null`)

			if err != nil {
				fmt.Fprintf(w, "%s\n", err)
				return
			}

			crns, err = json.Marshal(times.GroupByCrns())

			if err != nil {
				fmt.Fprintf(w, "%s\n", err)
				return
			}

			_, err = conn.Do("SET", "cache:crns", string(crns))

			if err != nil {
				fmt.Fprintf(w, "%s\n", err)
				return
			}

			conn.Do("EXPIRE", "cache:crns", "60")

		}

		w.Write(crns)
	})

	http.HandleFunc("/api/attrs", func(w http.ResponseWriter, r *http.Request) {

		conn := pool.Get()
		defer conn.Close()

		crns, err := redis.Bytes(conn.Do("GET", "cache:attrs"))

		if err != nil {
			attrs := []struct{
				Short string `json:"short"`
				Name string `json:"name"`
			}{}

			err := db.Select(&attrs, `select
	short,
	name
from
	course_attrs`)

			if err != nil {
				fmt.Fprintf(w, "%s\n", err)
				return
			}

			crns, err = json.Marshal(attrs)

			if err != nil {
				fmt.Fprintf(w, "%s\n", err)
				return
			}

			_, err = conn.Do("SET", "cache:attrs", string(crns))

			if err != nil {
				fmt.Fprintf(w, "%s\n", err)
				return
			}

			conn.Do("EXPIRE", "cache:attrs", "60")

		}

		w.Write(crns)
	})

	http.HandleFunc("/api/schedule", func(w http.ResponseWriter, r *http.Request) {


		startTime := time.Now()

		titles := r.URL.Query()["classTitles[]"]
		for i, title := range titles {
			titles[i] = fmt.Sprintf("title:%s", title)
		}

		attrs := r.URL.Query()["attrs[]"]
		for _, attr := range attrs {
			titles = append(titles, fmt.Sprintf("attr:%s", attr))
		}

		jsonBytes, err := json.Marshal(titles)
		if err != nil {
			fmt.Fprintf(w, "%s\n", err)
			return
		}

		conn, err := lentil.Dial(os.ExpandEnv("${BEANSTALKD_HOST}:${BEANSTALKD_PORT}"))
		if err != nil {
			fmt.Fprintf(w, "%s\n", err)
			return
		}
		conn.Use("schedules")
		jobId, err := conn.Put(0, 0, 60, jsonBytes)
		if err != nil {
			fmt.Fprintf(w, "%s\n", err)
			return
		}

		conn.Watch(fmt.Sprintf("schedule_result_%d", jobId))
		result, err := conn.Reserve()
		if err != nil {
			fmt.Fprintf(w, "%s\n", err)
			return
		}
		conn.Delete(result.Id)
		conn.Ignore(fmt.Sprintf("schedule_result_%d", jobId))

		log.Println(time.Since(startTime))

		w.Write(result.Body)


	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}