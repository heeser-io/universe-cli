FROM golang:1.18.3-alpine3.16

RUN apk add --no-cache ca-certificates git

RUN apk add --no-cache ca-certificates zip

COPY ./bin/linux64/universe /usr/local/bin/universe
# COPY crt /crt
# COPY .env /.env
