#!/bin/bash

POSTGRES="gosu postgres postgres"

$POSTGRES --single -E <<EOF
CREATE USER tally WITH PASSWORD 'tally' SUPERUSER;
CREATE DATABASE tally;
EOF

$POSTGRES --single -j -E tally < /app/db.sql