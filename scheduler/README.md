Scheduler
=========
The part that assembles the base requirements into all possible schedules.  Eventually this will be called/used by the job queue is a process like:

```
api -----> queue -----> scheduler
 ^          | ^             v
 |          v |             |
 ------------ ---------------
```

Required Environment Vars:
-------------------------
- BEANSTALKD_HOST
- BEANSTALKD_PORT
- REDIS_HOST
- REDIS_PORT