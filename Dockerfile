FROM alpine:3.18

LABEL org.opencontainers.image.source https://github.com/heeser-io/universe-cli

RUN apk add --no-cache ca-certificates git zip

COPY ./bin/linux64/universe /usr/local/bin/universe
# COPY crt /crt
# COPY .env /.env
