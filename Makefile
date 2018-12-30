GO=go
GOBIN=$(PWD)/bin
TESTOPTS=-coverprofile=result.coverprofile -v -race ./...
BUILDOPTS=-tags netgo -installsuffix netgo -ldflags "-w -s -extldflags -static"
BINARY=my-scraper

all: deps test build

setup:
	GOBIN=$(GOBIN) GO111MODULE=on go install github.com/lestrrat/go-server-starter/cmd/start_server
	GOBIN=$(GOBIN) GO111MODULE=on go install github.com/mattn/goveralls
	GOBIN=$(GOBIN) GO111MODULE=on go install honnef.co/go/tools/cmd/megacheck

deps:
	GO111MODULE=on go mod tidy

test:
	GO111MODULE=on $(GO) mod verify
	GO111MODULE=on $(GO) vet ./...
	GO111MODULE=on $(GO) test $(TESTOPTS)
	#$(GOBIN)/megacheck ./...

build:
	GO111MODULE=on $(GO) build -o $(BINARY) $(BUILDOPTS)

goveralls:
	$(GOBIN)/goveralls -service=travis-ci -coverprofile=result.coverprofile

build-linux:
	GO111MODULE=on CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GO) build -o $(BINARY).linux $(BUILDOPTS)
