FROM golang:1.7-alpine

RUN apk update && apk add git
ADD . /go/src/github.com/danohuiginn/grepurl
WORKDIR /go/src/github.com/danohuiginn/grepurl
RUN go get

