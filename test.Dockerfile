FROM docker.io/library/golang:1.16-buster as builder
LABEL maintainer="donato@wolfisberg.dev"
WORKDIR /app

COPY ./ .
RUN go build

FROM python:3.11.4-buster
WORKDIR /app

COPY ./.env ./server-updater-config.env
COPY --from=builder /app/server-updater ./
CMD ["./server-updater", "--env-file", "server-updater-config.env"]
