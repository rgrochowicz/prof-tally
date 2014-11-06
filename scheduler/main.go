package main

import (
	"encoding/json"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"log"
	"time"
)

var (
	c redis.Conn
)

type CrnSchedule struct {
	Crns []string
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

	//connect to redis
	c, err = redis.Dial("tcp", ":6379")
	if err != nil {
		log.Fatalln(err)
	}
	defer c.Close()

	titles := []string{"FOUNDATIONS COMP SCI", "OBJ-ORIENT PRGRM/DATA ABSTR", "LINEAR ALGEBRA", "COLLEGE COMPOSITION II", "PUBLIC SPEAKING", "DISCRETE STRUCTURES"}

	startTime := time.Now()

	crnSchedules := getNext([]string{}, titles)

	jsonBytes, err := json.Marshal(crnSchedules)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println(time.Since(startTime))

	log.Println(len(crnSchedules))
	log.Println(len(jsonBytes))

}
