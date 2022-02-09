ARG GO_VERSION=1.17

FROM golang:${GO_VERSION}-alpine AS builder

RUN mkdir /user \
    && echo 'daemon:x:2:2:daemon:/:' > /user/passwd \
    && echo 'daemon:x:2:' > /user/group

RUN apk add --no-cache ca-certificates git

RUN go get -u github.com/golang/dep/cmd/dep

WORKDIR ${GOPATH}/src/mtjobrunner

COPY ./Gopkg.toml ./Gopkg.lock ./
RUN dep ensure -vendor-only

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-s -w \
       -a -installsuffix 'static' -o /mtjobrunner ./cmd/*

FROM scratch AS mtjobrunner

COPY --from=builder /user/group /user/passwd /etc/
COPY --from=builder /mtjobrunner /mtjobrunner

USER daemon:daemon

ENTRYPOINT [ "/mtjobrunner" ]