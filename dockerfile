FROM golang:alpine

WORKDIR /go/src/github.com/jackyczj/NoGhost

ADD . /go/src/github.com/jackyczj/NoGhost

ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.cn

RUN go get && go build

EXPOSE 2333

ENTRYPOINT ./NoGhost