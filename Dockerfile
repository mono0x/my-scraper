FROM golang:1.11 AS builder

WORKDIR /go/src/github.com/mono0x/my-scraper

ADD . .
RUN make build-linux

FROM alpine:latest
RUN apk add --no-cache ca-certificates tzdata

COPY --from=builder /go/src/github.com/mono0x/my-scraper/my-scraper.linux /app
EXPOSE 8080
CMD ["/app"]
