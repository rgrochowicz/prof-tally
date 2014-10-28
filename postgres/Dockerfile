FROM postgres:9.3

RUN mkdir /docker-entrypoint-initdb.d

ADD ./db.sql /app/db.sql
ADD ./init.sh /docker-entrypoint-initdb.d/init.sh