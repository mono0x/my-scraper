GO=go
GOBIN=$(PWD)/bin
TESTOPTS=-coverprofile=result.coverprofile -v -race ./...
BUILDOPTS=-tags netgo -installsuffix netgo -ldflags "-w -s -extldflags -static"
BINARY=my-scraper

all: deps test build

setup:
	GOBIN=$(GOBIN) GO111MODULE=on $(GO) install github.com/lestrrat-go/server-starter/cmd/start_server
	GOBIN=$(GOBIN) GO111MODULE=on $(GO) install honnef.co/go/tools/cmd/staticcheck

download:
	GO111MODULE=on $(GO) mod download

deps:
	GO111MODULE=on $(GO) mod tidy

upgrade-deps:
	GO111MODULE=on $(GO) get -u
	GO111MODULE=on $(GO) mod tidy

test:
	GO111MODULE=on $(GO) mod verify
	GO111MODULE=on $(GO) vet ./...
	GO111MODULE=on $(GO) test $(TESTOPTS)
	GO111MODULE=on $(GOBIN)/staticcheck ./...

build:
	GO111MODULE=on $(GO) build -o $(BINARY) $(BUILDOPTS)

build-linux:
	GO111MODULE=on CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GO) build -o $(BINARY).linux $(BUILDOPTS)
