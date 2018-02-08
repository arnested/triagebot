FROM golang:latest AS build-env

WORKDIR /go/src/github.com/arnested/triagebot
COPY . .

RUN go version
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o triagebot .

FROM scratch

COPY --from=build-env /etc/ssl/certs/ /etc/ssl/certs/
COPY --from=build-env /go/src/github.com/arnested/triagebot/triagebot /triagebot

ENTRYPOINT ["/triagebot"]
