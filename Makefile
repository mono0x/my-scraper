GO=go
TESTOPTS=-coverprofile=result.coverprofile -v -race ./...
BUILDOPTS=-tags netgo -installsuffix netgo -ldflags "-w -s -extldflags -static"
BINARY=my-scraper

all: deps test build

setup:
	go get -u github.com/twitchtv/retool

deps:
	retool sync
	retool do dep ensure

test:
	retool do megacheck ./...
	$(GO) vet ./...
	$(GO) test $(TESTOPTS)

build:
	$(GO) build -o $(BINARY) $(BUILDOPTS)

build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GO) build -o $(BINARY).linux $(BUILDOPTS)
