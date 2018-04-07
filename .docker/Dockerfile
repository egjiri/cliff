FROM golang:1.10.0-alpine3.7 as builder

RUN apk update && apk add git

RUN go get -u github.com/jteeuwen/go-bindata/...

WORKDIR /go/src/github.com/egjiri/cliff

COPY . .

RUN go run .docker/main.go

VOLUME /data

CMD go-bindata -pkg data -o data/go-bindata.go -prefix=/data/ /data/cli.yml && go build -o /data/cliff

LABEL image="egjiri/cliff:0.0.1"
