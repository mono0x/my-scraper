FROM ghcr.io/jdx/mise:2026.6.1@sha256:b319ba78c02fd28a593d043dd09d62685191420f327bf1534bdaa128f896bc93 AS builder

ENV GOPATH=/go
WORKDIR /go/src/github.com/mono0x/my-scraper

RUN --mount=type=bind,source=mise.toml,target=mise.toml \
    mise trust && mise install

RUN --mount=type=bind,source=mise.toml,target=mise.toml \
    --mount=type=bind,source=go.mod,target=go.mod \
    --mount=type=bind,source=go.sum,target=go.sum \
    --mount=type=cache,target=/go/pkg/mod,sharing=locked \
    go mod download

RUN --mount=type=bind,source=.,target=. \
    --mount=type=cache,target=/go/pkg/mod \
    go tool task build -- -o /bin/my-scraper

FROM gcr.io/distroless/static-debian12:nonroot@sha256:d093aa3e30dbadd3efe1310db061a14da60299baff8450a17fe0ccc514a16639

COPY --from=builder --chown=nonroot:nonroot /bin/my-scraper /bin/my-scraper
CMD ["/bin/my-scraper"]
