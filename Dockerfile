FROM golang:1.17-alpine as dev
RUN apk add build-base && apk add sqlite && go get github.com/mattn/go-sqlite3@v1.14.10
WORKDIR /work
