Populate
========

Populates redis with data from postgres.

Required Environment Vars:
--------------------------
- POSTGRES_DATABASE
- POSTGRES_USER
- POSTGRES_PASSWORD
- POSTGRES_HOST
- POSTGRES_PORT
- REDIS_HOST
- REDIS_PORT

How population works:
---------------------

- For every crn in the database, a set is made on redis that contains all crns that don't overlap that crn.
- For every title, a set of crns that have that title.
- For every attribute, a set of crns that have that attribute.