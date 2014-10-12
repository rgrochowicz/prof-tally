Scraper
=======

This is the scraper portion. Scrape the section tally page and dump it into a database.

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