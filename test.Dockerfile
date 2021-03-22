FROM docker.io/library/golang:1.16-buster as builder
LABEL maintainer="donato@wolfisberg.dev"
WORKDIR /app

COPY ./ .
RUN go build
RUN ./server-updater --env-file server-updater-config.env
