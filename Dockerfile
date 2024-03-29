FROM golang:1.22-bookworm AS builder

WORKDIR /go/src/github.com/mono0x/my-scraper

COPY go.mod go.sum Makefile ./
RUN make download

COPY . ./
RUN make build-linux

# hadolint ignore=DL3006
FROM gcr.io/distroless/static-debian12

COPY --from=builder /go/src/github.com/mono0x/my-scraper/my-scraper.linux /app
CMD ["/app"]
