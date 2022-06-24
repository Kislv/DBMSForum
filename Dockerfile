FROM golang:1.15.2-buster AS build

RUN mkdir /app

COPY . /app

WORKDIR /app

RUN go build -o api-server main.go

FROM ubuntu:20.04 AS release

# Make the "en_US.UTF-8" locale so postgres will be utf-8 enabled by default
RUN apt-get -y update && apt-get install -y locales gnupg2
RUN locale-gen en_US.UTF-8
RUN update-locale LANG=en_US.UTF-8

# install Postgres
ENV PGVER 12
ENV DEBIAN_FRONTEND noninteractive
RUN apt-get update -y && apt-get install -y postgresql postgresql-contrib

# Run the rest of the commands as the ``postgres`` user created by the ``postgres-$PGVER`` package when it was ``apt installed``
USER postgres

# Create a PostgreSQL role named ``vitya`` with ``password`` as the password and
# then create a database `forum` owned by the ``vitya`` role.
RUN /etc/init.d/postgresql start &&\
    psql --command "CREATE USER vitya WITH SUPERUSER PASSWORD 'password';" &&\
    createdb -E UTF8 forum &&\
    /etc/init.d/postgresql stop

RUN echo "listen_addresses='*'\n" >> /etc/postgresql/$PGVER/main/postgresql.conf
RUN echo "host all  all    0.0.0.0/0  md5" >> /etc/postgresql/$PGVER/main/pg_hba.conf

# Expose the PostgreSQL port

# Add VOLUMEs to allow backup of config, logs and databases
VOLUME ["/etc/postgresql", "/var/log/postgresql", "/var/lib/postgresql"]

# Back to the root user
USER root

# Собранный сервер
COPY --from=build /app /app

EXPOSE 5432
EXPOSE 5000

ENV PGPASSWORD password

CMD service postgresql start && psql -h localhost -d forum -U vitya -p 5432 -a -q -f /app/script/init.sql && /app/api-server