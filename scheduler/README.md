Scheduler
=========
The part that assembles the base requirements into all possible schedules.  Eventually this will be called/used by the job queue is a process like:

```
api -----> beanstalkd -----> scheduler
 ^          |      ^             v
 |          v      |             |
 ------------      ---------------
```

Required Environment Vars:
-------------------------
- BEANSTALKD_HOST
- BEANSTALKD_PORT
- REDIS_HOST
- REDIS_PORT

How this works:
---------------

1. The scheduler receives a job from beanstalkd on the `schedules` pipe
2. Needles are parsed from the job data (eg. `title:COMP II`, `attr:RSEM`)
3. The crns for the first needle are retrieved from redis.
4. For every crn retrieved, the crns that intersect the crns that don't overlap that crn and crns that are in the next needle are retrieved.
5. This process continues until the needle list is exhausted.
6. The result is json-encoded and sent back through beanstalkd to a reply channel named `'schedule_result_' + the job's id`

Example:
--------

```
needles: ["title:A", "title:B"] 
# we want all schedules that contain these two needles

get the crns for 'title:A' -> 1,2,3,4
loop through those crns

for the first crn:
	intersect (the crns that don't overlap crn:1) and (the crns in title:B)
	add the result of the intersection to the final results

for the second crn:
	intersect (the crns that don't overlap crn:2) and (the crns in title:B)
	add the result of the intersection to the final results

```