# syntax=docker/dockerfile:1
# check=error=true

FROM ghcr.io/jdx/mise:2026.6.14@sha256:b8f8c20fc3308f8b1d00ccca2bc968e4e208af1c5c1069e1ad9753baa099acff AS builder

ENV GOPATH=/go
WORKDIR /go/src/github.com/mono0x/my-scraper

RUN --mount=type=bind,source=mise.toml,target=mise.toml \
    mise trust && mise install 

RUN --mount=type=bind,source=mise.toml,target=mise.toml \
    --mount=type=bind,source=go.mod,target=go.mod \
    --mount=type=bind,source=go.sum,target=go.sum \
    --mount=type=cache,target=/go/pkg/mod,sharing=locked \
    mise exec -- go mod download

RUN --mount=type=bind,source=.,target=. \
    --mount=type=cache,target=/go/pkg/mod \
    mise exec -- go tool task build -- -o /bin/my-scraper

FROM gcr.io/distroless/static-debian12:nonroot@sha256:d093aa3e30dbadd3efe1310db061a14da60299baff8450a17fe0ccc514a16639

COPY --from=builder --chown=nonroot:nonroot /bin/my-scraper /bin/my-scraper
CMD ["/bin/my-scraper"]
