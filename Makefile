VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
LDFLAGS := -ldflags "-X main.version=$(VERSION)"
BINARY := kidmode

.PHONY: build clean

build:
	go build $(LDFLAGS) -o $(BINARY) .

clean:
	rm -f $(BINARY)
