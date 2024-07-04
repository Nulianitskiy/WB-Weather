FROM postgres:alpine3.18

ENV POSTGRES_PASSWORD=caravanpalace \
    POSTGRES_USER=wbuser \
    POSTGRES_DB=wbweather

COPY ./db.sql /docker-entrypoint-initdb.d/db.sql

EXPOSE 5432