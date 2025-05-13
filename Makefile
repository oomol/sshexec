# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
BINARY_NAME=sshexec

# All target
all: build

build: build-arm64 build-amd64

build-arm64:
	GOARCH=arm64 $(GOBUILD) -o out/$(BINARY_NAME)-arm64 -v sshd/cmd/
	GOOS=linux GOARCH=arm64 $(GOBUILD) -o out/caller-arm64 -v scripts/caller.go


build-amd64:
	GOARCH=amd64 $(GOBUILD) -o out/$(BINARY_NAME)-amd64 -v  sshd/cmd/
	GOOS=linux GOARCH=amd64 $(GOBUILD) -o out/caller-amd64 -v scripts/caller.go

lint:
	golangci-lint run

# Clean the build files
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)

fmt:
	$(GOCMD) fmt ./...

.PHONY: all build clean fmt run
