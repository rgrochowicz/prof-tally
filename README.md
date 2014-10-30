Prof Tally
==========

This will eventually make a schedule automatically based on inputs.

It will comprise of three parts:

1. Scraper - Scrapes details from the tally page
2. API - Provides access from the web components
3. Web - The interactive part
4. Postgres - Gives a place to store data

Deployment
----------

All of the modules are available on the docker repositiory.

Run them in this order:

1. To run the postgres server, run `rgrochowicz/prof-tally:postgres`.  The database schema will be automatically applied using `postgres/db.sql`
2. To run the API endpoint, run `rgrochowicz/prof-tally:api` with the environment variables `POSTGRES_USER` `POSTGRES_PASSWORD` `POSTGRES_DATABASE` `POSTGRES_HOST` `POSTGRES_PORT`
3. To run the public-facing web service, run `rgrochowicz/prof-tally:web` with the environment variables `API_HOST` `API_PORT`
4. To populate the data, run `rgrochowicz/prof-tally:scraper` with the environment variables `POSTGRES_USER` `POSTGRES_PASSWORD` `POSTGRES_DATABASE` `POSTGRES_HOST` `POSTGRES_PORT`

The respective docker commands would be:

1. `docker run --rm -p 5432:5432 rgrochowicz/prof-tally:postgres`
2. `docker run --rm -p 3000:3000 -e POSTGRES_USER=tally -e POSTGRES_PASSWORD=tally -e POSTGRES_DATABASE=tally -e POSTGRES_HOST=127.0.0.1 -e POSTGRES_PORT=5432 rgrochowicz/prof-tally:api`
3. `docker run --rm -p 80:80 -e API_HOST=127.0.0.1 -e API_PORT=3000 rgrochowicz/prof-tally:web`
4. `docker run --rm -e POSTGRES_USER=tally -e POSTGRES_PASSWORD=tally -e POSTGRES_DATABASE=tally -e POSTGRES_HOST=127.0.0.1 -e POSTGRES_PORT=5432 rgrochowicz/prof-tally:scraper`

To use this with CoreOS, fleet unit files are provided.  Launch the services with:

1. `cd` into the `unit-files` directory
1. `fleetctl start postgres@1.service postgres-discovery@1.service`
2. `fleetctl start api@1.service api-discovery@1.service`
3. `fleetctl start web@1.service web-discovery@1.service`
4. `fleetctl start scraper@1.service`