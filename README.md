Prof Tally
==========

This will eventually make a schedule automatically based on inputs.

It will comprise of 8 parts:

1. Scraper - Scrapes details from the tally page
2. API - Provides access from the web components
3. Web - The interactive part
4. Populate - Computes data for redis
5. Scheduler - Produces a bunch of schedules
4. Postgres - Gives a place to store data
5. Redis - Data for schedule computation
7. Beanstalkd - Job queue for api <--> scheduler

Deployment
----------

All of the modules are available on the docker repositiory.

Run them in this order:

1. To run the postgres server, run `rgrochowicz/prof-tally:postgres`.  The database schema will be automatically applied using `postgres/db.sql`
2. To run the API endpoint, run `rgrochowicz/prof-tally:api`
3. To run the public-facing web service, run `rgrochowicz/prof-tally:web`
4. To populate the data, run `rgrochowicz/prof-tally:scraper`
5. To run redis, run `rgrochowicz/prof-tally:redis`
6. To run beanstalkd, run `rgrochowicz/prof-tally:beanstalkd`
7. To run populate, run `rgrochowicz/prof-tally:populate`
8. To run the scheduler, run `rgrochowicz/prof-tally:scheduler`

To use this with CoreOS, fleet unit files are provided.  Launch the services with:

1. `cd` into the `unit-files` directory
2. `fleetctl start postgres@1.service postgres-discovery@1.service`
3. `fleetctl start redis@1.service redis-discovery@1.service`
4. `fleetctl start beanstalkd@1.service beanstalkd-discovery@1.service`
5. `fleetctl start api@1.service api-discovery@1.service`
6. `fleetctl start web@1.service web-discovery@1.service`
7. `fleetctl start scraper@1.service`
8. `fleetctl start populate@1.service`
9. `fleetctl start scheduler@1.service`