# syntax=docker/dockerfile:1

# Build

FROM golang:1.19 AS build

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY *.go .

RUN go build -o /hello

# Deploy 

FROM debian:latest

WORKDIR /

COPY --from=build /hello /usr/local/bin/hello

EXPOSE 3333

ENTRYPOINT [ "/usr/local/bin/hello" ]