FROM golang:alpine

LABEL maintainer = "marcoding78@gmail.com"

# Set pwd to the go folder
WORKDIR /movieCatcherApp

ADD . .

RUN go mod download

ENTRYPOINT go build && ./RESTfulTest01