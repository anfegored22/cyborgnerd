.PHONY: dev test binary

BIN ?= ../bin/cyborgnerd

dev:
	go run .
test:
	go test ./...
binary:
	CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o "$(BIN)" .
