FROM golang:1.12-alpine AS build-env

RUN apk add git

ENV GO111MODULE=on

WORKDIR /build
COPY *.go go.mod go.sum /build/


ENV CGO_ENABLED=0
ENV GOOS=linux

RUN go version
RUN go build
RUN go test -o ./triagebot.test -v -cover ./...

FROM scratch

ENV PATH=/

COPY --from=build-env /usr/local/go/lib/time/zoneinfo.zip /usr/local/go/lib/time/zoneinfo.zip
COPY --from=build-env /etc/ssl/certs/ /etc/ssl/certs/
COPY --from=build-env /build/triagebot /triagebot
COPY --from=build-env /build/triagebot.test /test

ENTRYPOINT ["triagebot"]
