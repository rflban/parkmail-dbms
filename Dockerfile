FROM golang:latest AS builder
WORKDIR /bundle
COPY . ./
RUN GOAMD64=v3 go build github.com/rflban/parkmail-dbms/cmd/forum

FROM ubuntu:20.04
RUN apt-get -y update && apt-get install -y tzdata
RUN ln -snf /usr/share/zoneinfo/Russia/Moscow /etc/localtime && echo Russia/Moscow > /etc/timezone
RUN apt-get -y update && apt-get install -y postgresql-12 && rm -rf /var/lib/apt/lists/*
USER postgres
RUN /etc/init.d/postgresql start && psql --command "CREATE USER forum WITH SUPERUSER PASSWORD 'forum';" && createdb -O forum forum && /etc/init.d/postgresql stop
EXPOSE 5432
WORKDIR /app
COPY ./configs ./configs
COPY --from=builder /bundle/forum ./forum
EXPOSE 5000
ENV PGPASSWORD forum
CMD service postgresql start && psql -h localhost -d forum -U forum -p 5432 -a -q -f ./configs/sql/init.sql && ./forum
