package main

import (
	"encoding/json"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"github.com/nutrun/lentil"
	"log"
	"time"
)

var (
	c redis.Conn
)

type CrnSchedule []string

func getNext(crns []string, titles []string) []CrnSchedule {

	//when there are not titles left, return the crns collected
	if len(titles) == 0 {
		crnsCopy := make(CrnSchedule, len(crns))
		copy(crnsCopy, crns)

		return []CrnSchedule{
			crnsCopy,
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

func MakeSchedules(titles []string) ([]byte, error) {
	log.Println(titles)
	startTime := time.Now()

	crnSchedules := getNext([]string{}, titles)

	jsonBytes, err := json.Marshal(crnSchedules)
	if err != nil {
		return nil, err
	}

	log.Println(time.Since(startTime))

	return jsonBytes, nil
}

func main() {

	var err error

	//connect to redis
	c, err = redis.Dial("tcp", ":6379")
	if err != nil {
		log.Fatalln(err)
	}
	defer c.Close()

	// Limit number of concurrent jobs execution. Use worker.Unlimited (0) if you want no limitation.
	conn, err := lentil.Dial("0.0.0.0:11300")
	if err != nil {
		log.Fatalln(err)
	}

	conn.Watch("schedules")

	log.Println("Listening")
	for {
		job, err := conn.Reserve()
		if err != nil {
			log.Fatalln(err)
		}

		titles := []string{}
		err = json.Unmarshal(job.Body, &titles)
		if err != nil {
			log.Fatalln(err)
		}

		result, err := MakeSchedules(titles)
		if err != nil {
			log.Fatalln(err)
		}

		err = conn.Delete(job.Id)

		if err != nil {
			log.Fatalln(err)
		}

		conn.Use(fmt.Sprintf("schedule_result_%d", job.Id))
		conn.Put(0, 0, 60, result)

	}

	log.Println("Exiting")

}
