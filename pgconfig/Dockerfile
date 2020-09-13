FROM postgres
COPY ./pg.sql docker-entrypoint-initdb.d/pg.sql
ENV POSTGRES_PASSWORD distributed
ENV POSTGRES_USER distributed
ENV POSTGRES_DB distributed