ARG GO_VERSION=1.17

FROM golang:${GO_VERSION}-alpine AS builder

RUN mkdir /user \
    && echo 'daemon:x:2:2:daemon:/:' > /user/passwd \
    && echo 'daemon:x:2:' > /user/group

RUN apk add --no-cache ca-certificates git gcc musl-dev librdkafka-dev pkgconf

WORKDIR ${GOPATH}/src/mtjobrunner

COPY client ./client
COPY cmd ./cmd
COPY pkg ./pkg
COPY go.mod ./go.mod
COPY go.sum ./go.sum

#RUN go get k8s.io/client-go@v0.23.3 \
#    github.com/confluentinc/confluent-kafka-go/kafka@v1.8.2 \
#    go.uber.org/zap@v1.21.0

RUN GOOS=linux go build -tags musl -tags dynamic -ldflags "-s -w" \
       -a -installsuffix 'static' -o /mtjobrunner ./cmd/*

RUN go clean --modcache
RUN apk del gcc musl-dev librdkafka-dev pkgconf
#RUN apk cache clean

USER daemon:daemon

ENTRYPOINT [ "/mtjobrunner" ]