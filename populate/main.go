package main;

import (
	_ "github.com/lib/pq"
    _ "database/sql"
    "github.com/jmoiron/sqlx"
    "log"
    "github.com/garyburd/redigo/redis"
    "fmt"
    "os"
)

var (
	db *sqlx.DB
	c redis.Conn
)

func clear() {

	//remove all contents from redis
	c.Do("FLUSHDB")
}

func crns() {

	//get all crns from postgres
    var crns []int
    err := db.Select(&crns, "select crn from courses")
    if err != nil {
    	log.Fatalln(err)
    }

    //for each crn, get all crns that don't overlap in times
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

    	//add the non-overlapping crns to a set named 'crn:' + crn
		c.Send("SADD", redis.Args{}.Add(fmt.Sprintf("crn:%d", crn)).AddFlat(nonoverlappingCrns)...)
    }

    //send the commands to redis
    c.Flush()

    //get back its replies
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

	//send the commands to redis
	c.Flush()

	//get back its replies
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

	//send the commands to redis
	c.Flush()

	//get back its replies
	for range contents {
		c.Receive()
	}

}

func main() {

	var err error

	//connect to postgres
	db, err = sqlx.Connect("postgres", os.ExpandEnv("user=${POSTGRES_USER} dbname=${POSTGRES_DATABASE} password=${POSTGRES_PASSWORD} host=${POSTGRES_HOST} port=${POSTGRES_PORT} sslmode=disable"))
    if err != nil {
        log.Fatalln(err)
    }
    defer db.Close()

    //connect to redis
    c, err = redis.Dial("tcp", os.ExpandEnv("${REDIS_HOST}:${REDIS_PORT}"))
    if err != nil {
    	log.Fatalln(err)
    }
    defer c.Close()

    //clear out redis
    clear()

    //insert all crns
    crns()

    //insert all titles
    titles()

    //insert all attributes
    attributes()

}