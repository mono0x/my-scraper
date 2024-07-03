FROM golang:1.22.5-bookworm AS builder

WORKDIR /go/src/github.com/mono0x/my-scraper

COPY go.mod go.sum Makefile ./
RUN make download

COPY . ./
RUN make build-linux

FROM gcr.io/distroless/static-debian12:nonroot

COPY --from=builder --chown=nonroot:nonroot /go/src/github.com/mono0x/my-scraper/my-scraper.linux /app
CMD ["/app"]
