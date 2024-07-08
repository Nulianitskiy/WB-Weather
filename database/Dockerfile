FROM postgres:alpine3.18

ENV POSTGRES_PASSWORD=caravanpalace \
    POSTGRES_USER=wbuser \
    POSTGRES_DB=wbweather

COPY migrate.sql /docker-entrypoint-initdb.d/migrate.sql

EXPOSE 5432