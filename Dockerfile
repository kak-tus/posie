FROM golang:1.12.9-alpine3.10 AS build

WORKDIR /go/posie

ENV CGO_ENABLED=0

COPY go.mod .
COPY go.sum .
COPY main.go .
COPY vendor ./vendor

RUN go test -mod=vendor && go build -mod=vendor -o /go/bin/posie

FROM alpine:3.10

COPY --from=build /go/bin/posie /usr/local/bin/posie

RUN apk add --no-cache \
    ca-certificates

ENV \
  POSIE_ADDR= \
  POSIE_TEXT_FAIL= \
  POSIE_TEXT_OK= \
  POSIE_TG_CHAT= \
  POSIE_TG_TOKEN=

CMD ["/usr/local/bin/posie"]
