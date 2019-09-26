FROM golang:alpine

WORKDIR /go/src/github.com/jackyczj/July
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories && \
        apk add gcc musl-dev
ADD . /go/src/github.com/jackyczj/July

ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.cn

RUN go get && go build

EXPOSE 2333

ENTRYPOINT ./July