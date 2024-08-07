FROM golang:1.22.5-bookworm AS builder

WORKDIR /go/src/github.com/mono0x/my-scraper

RUN --mount=type=bind,source=go.mod,target=go.mod \
    --mount=type=bind,source=go.sum,target=go.sum \
    --mount=type=bind,source=Makefile,target=Makefile \
    --mount=type=cache,target=/go/pkg/mod,sharing=locked \
    make download

RUN --mount=type=bind,source=.,target=. \
    --mount=type=cache,target=/go/pkg/mod \
    make build-docker

FROM gcr.io/distroless/static-debian12:nonroot

COPY --from=builder --chown=nonroot:nonroot /bin/my-scraper /bin/my-scraper
CMD ["/bin/my-scraper"]
