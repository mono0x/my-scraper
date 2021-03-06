FROM golang:1.16 AS builder

WORKDIR /go/src/github.com/mono0x/my-scraper

ADD go.mod go.sum Makefile ./
RUN make download

ADD . ./
RUN make build-linux

FROM scratch
COPY --from=golang:1.16 /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

COPY --from=builder /go/src/github.com/mono0x/my-scraper/my-scraper.linux /app
CMD ["/app"]
