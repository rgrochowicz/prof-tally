package main

import (
	"encoding/json"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"github.com/nutrun/lentil"
	"log"
	"time"
	"os"
)

type CrnSchedule []string

type Scheduler struct {
	Conn redis.Conn
	ScheduleResult []CrnSchedule
}

func NewScheduler() (*Scheduler, error) {

	//connect to redis
	c, err := redis.Dial("tcp", os.ExpandEnv("${REDIS_HOST}:${REDIS_PORT}"))
	if err != nil {
		return nil, err
	}

	return &Scheduler{
		Conn: c,
	}, nil
}

func (s *Scheduler) Make(needles []string) []CrnSchedule {
	
	//clear out the previous results
	s.ScheduleResult = []CrnSchedule{}

	s.getNext([]string{}, needles)

	return s.ScheduleResult
}

func (s *Scheduler) getNext(crns, needles []string) {

	nextNeedle := needles[0]

	crnArgs := []string{}
	for _, crn := range crns {
		crnArgs = append(crnArgs, fmt.Sprintf("crn:%s", crn))
	}

	//get all available courses based on set intersections
	entries, err := redis.Strings(s.Conn.Do("SINTER", redis.Args{}.AddFlat(crnArgs).Add(nextNeedle)...))
	if err != nil {
		log.Fatalln(err)
	}

	//when we're on the last needle, append the crns collected
	if len(needles) == 1 {
		for _, entry := range entries {

			crnsCopy := make(CrnSchedule, len(crns))
			copy(crnsCopy, crns)
			crnsCopy = append(crnsCopy, entry)

			s.ScheduleResult = append(s.ScheduleResult, crnsCopy)
		}

		return
	}

	for _, entry := range entries {
		s.getNext(append(crns, entry), needles[1:])
	}
}

func (s *Scheduler) Close() {
	s.Conn.Close()
}



func MakeSchedules(needles []string) ([]byte, error) {
	log.Println(needles)
	startTime := time.Now()

	sched, err := NewScheduler()
	if err != nil {
		log.Fatalln(err)
	}
	defer sched.Close()

	crnSchedules := sched.Make(needles)

	log.Printf("Schedules made: %d", len(crnSchedules))

	jsonTime := time.Now()
	jsonBytes, err := json.Marshal(crnSchedules)
	if err != nil {
		return nil, err
	}

	log.Printf("JSON time: %s", time.Since(jsonTime))
	log.Printf("Total time: %s", time.Since(startTime))

	return jsonBytes, nil
}

func main() {

	//connect to beanstalkd
	conn, err := lentil.Dial(os.ExpandEnv("${BEANSTALKD_HOST}:${BEANSTALKD_PORT}"))
	if err != nil {
		log.Fatalln(err)
	}

	conn.Watch("schedules")

	log.Println("Listening: scheduler")
	for {

		//get a job from beanstalkd
		job, err := conn.Reserve()
		if err != nil {
			log.Fatalln(err)
		}

		//get the needles for that job
		needles := []string{}
		err = json.Unmarshal(job.Body, &needles)
		if err != nil {
			log.Fatalln(err)
		}

		//make the schedule based on those needles
		result, err := MakeSchedules(needles)
		if err != nil {
			log.Fatalln(err)
		}

		//delete the job from beanstalkd as it's complete
		err = conn.Delete(job.Id)

		if err != nil {
			log.Fatalln(err)
		}

		//send the result down the result tube
		conn.Use(fmt.Sprintf("schedule_result_%d", job.Id))
		conn.Put(0, 0, 60, result)

	}

	log.Println("Exiting")

}
