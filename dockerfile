FROM golang:last

WORKDIR /go/src/github.com/jackyczj/NoGhost

ADD . /go/src/github.com/jackyczj/NoGhost

RUN go get && go build

EXPOSE 2333

ENTRYPOINT ./NoGhost