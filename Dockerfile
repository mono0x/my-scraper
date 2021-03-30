FROM golang:1.16-buster AS builder

WORKDIR /go/src/github.com/mono0x/my-scraper

ADD go.mod go.sum Makefile ./
RUN make download

ADD . ./
RUN make build-linux

FROM gcr.io/distroless/static-debian10

COPY --from=builder /go/src/github.com/mono0x/my-scraper/my-scraper.linux /app
CMD ["/app"]
