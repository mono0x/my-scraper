GO=go
GOBIN=$(PWD)/bin
TESTOPTS=-v -race ./...
BUILDOPTS=-tags netgo,timetzdata -installsuffix netgo -ldflags "-w -s -extldflags -static"
BINARY=my-scraper

all: deps test build

setup:
	GOBIN=$(GOBIN) $(GO) install github.com/cosmtrek/air
	GOBIN=$(GOBIN) $(GO) install honnef.co/go/tools/cmd/staticcheck

download:
	$(GO) mod download

deps:
	$(GO) mod tidy

upgrade-deps:
	$(GO) get -u
	$(GO) mod tidy

test:
	$(GO) mod verify
	$(GO) vet ./...
	$(GO) test $(TESTOPTS)
	$(GOBIN)/staticcheck ./...

build:
	$(GO) build -o $(BINARY) $(BUILDOPTS)

build-docker:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GO) build -o /bin/$(BINARY) $(BUILDOPTS)

watch:
	$(GOBIN)/air
