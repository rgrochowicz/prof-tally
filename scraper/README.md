Scraper
=======

This is the scraper portion. Scrapes the section tally page and dumps it into a database.

Required Environment Vars:
--------------------------
- POSTGRES_DATABASE
- POSTGRES_USER
- POSTGRES_PASSWORD
- POSTGRES_HOST
- POSTGRES_PORT

Prereqs:
--------
1. Install postgres with a user named `tally`
2. Make postgres listen on localhost or *

How to get started:
-------------------
1. Install requirements using pip: `pip install -r requirements.txt`
2. Make database schema using db.sql

Running:
--------
1. Run with `python scrape.py`
2. Enjoy results


Limitations:
------------
1. Some classes have invalid times