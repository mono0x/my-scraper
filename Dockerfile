FROM ghcr.io/jdx/mise:2026.5.5@sha256:43f0895b1c3d4784f8ac2e5ff15ede056842570849ca034f2ecfe7dbdc005cc8 AS builder

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
