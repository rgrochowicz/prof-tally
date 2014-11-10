package main;

import (
	_ "github.com/lib/pq"
    _ "database/sql"
    "github.com/jmoiron/sqlx"
    "log"
    "github.com/garyburd/redigo/redis"
    "fmt"
)

var (
	db *sqlx.DB
	c redis.Conn
)

func clear() {
	c.Do("FLUSHDB")
}

func crns() {
    var crns []int
    err := db.Select(&crns, "select crn from courses")
    if err != nil {
    	log.Fatalln(err)
    }

    for _, crn := range crns {

    	var nonoverlappingCrns []string
    	err = db.Select(&nonoverlappingCrns, `with taken_times as (
	select weekday, start_time, length
	from course_times
	where course_crn = $1
), excluded_crns as (
	select course_crn
	from course_times
	where exists (
		select 1
		from taken_times
		where 
			(taken_times.start_time, taken_times.length) overlaps (course_times.start_time, course_times.length) and
			weekday = course_times.weekday
	)
)
select distinct(course_crn)
from course_times
where course_crn not in (select course_crn from excluded_crns)`, crn)
    	if err != nil {
    		log.Fatalln(err)
    	}

		c.Send("SADD", redis.Args{}.Add(fmt.Sprintf("crn:%d", crn)).AddFlat(nonoverlappingCrns)...)
    }

    c.Flush()

    for range crns {
    	c.Receive()
    }
}

func titles() {
	contents := []struct {
		Crn string
		Title string
	}{}

	//only get courses that have times
	db.Select(&contents, `select
	crn,
	title
from courses
where exists(select 1 from course_times where course_crn = crn)`)

	for _, content := range contents {
		c.Send("SADD", fmt.Sprintf("title:%s", content.Title), content.Crn)
	}
	c.Flush()
	for range contents {
		c.Receive()
	}
}

func attributes() {
	contents := []struct {
		Crn string
		Attr string
	}{}

	//only get courses that have times
	db.Select(&contents, `select
	crn,
	attr
from courses_and_attrs
where exists(select 1 from course_times where course_crn = crn)`)

	for _, content := range contents {
		c.Send("SADD", fmt.Sprintf("attr:%s", content.Attr), content.Crn)
	}
	c.Flush()
	for range contents {
		c.Receive()
	}

}

func main() {

	var err error

	db, err = sqlx.Connect("postgres", "user=tally dbname=tally password=tally sslmode=disable")
    if err != nil {
        log.Fatalln(err)
    }
    defer db.Close()

    c, err = redis.Dial("tcp", ":6380")
    if err != nil {
    	log.Fatalln(err)
    }
    defer c.Close()

    clear()
    crns()
    titles()
    attributes()

}