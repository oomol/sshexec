# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
BINARY_NAME=sshexec

# All target
all: build

build: build-arm64 build-amd64 script

build-arm64:
	GOARCH=arm64 $(GOBUILD) -o out/$(BINARY_NAME)-arm64 -v sshd/cmd/

build-amd64:
	GOARCH=amd64 $(GOBUILD) -o out/$(BINARY_NAME)-amd64 -v  sshd/cmd/

lint:
	golangci-lint run

# Clean the build files
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)

# Format the code
script:
	cp scripts/caller.sh out/
	chmod +x out/caller.sh
	cp scripts/installer.sh out/
	chmod +x out/installer.sh

fmt:
	$(GOCMD) fmt ./...

.PHONY: all build clean fmt run
