FROM golang:1.11 AS builder

ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

WORKDIR /go/src/github.com/mono0x/my-scraper

RUN go get github.com/golang/dep/cmd/dep

ADD Gopkg.toml Gopkg.lock ./
RUN /go/bin/dep ensure -vendor-only

ADD . .
RUN /go/bin/dep ensure && go build -tags netgo -installsuffix netgo -ldflags "-w -s -extldflags -static"

FROM alpine:latest
RUN apk add --no-cache ca-certificates tzdata

COPY --from=builder /go/src/github.com/mono0x/my-scraper/my-scraper /app
EXPOSE 8080
CMD ["/app"]
