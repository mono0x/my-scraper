# syntax=docker/dockerfile:1
# check=error=true

FROM ghcr.io/jdx/mise:2026.8.2@sha256:010e6829a39dafea7a660d15d72b730023fb36ab27c65ed0d65c325d6c979a61 AS builder

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
    mise run build -o /bin/my-scraper

FROM gcr.io/distroless/static-debian12:nonroot@sha256:afa5c872c891853ca7fcf1f12c3edb23f7eeef36189728842dd51042ff57f7ab

COPY --from=builder --chown=nonroot:nonroot /bin/my-scraper /bin/my-scraper
CMD ["/bin/my-scraper"]
