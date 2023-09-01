FROM golang:1.21-bullseye AS builder

WORKDIR /go/src/github.com/mono0x/my-scraper

COPY go.mod go.sum Makefile ./
RUN make download

COPY . ./
RUN make build-linux

# hadolint ignore=DL3006
FROM gcr.io/distroless/static-debian11

COPY --from=builder /go/src/github.com/mono0x/my-scraper/my-scraper.linux /app
CMD ["/app"]
