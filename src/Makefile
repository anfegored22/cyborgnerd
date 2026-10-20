.PHONY: tw dev css binary build test

BIN ?= ../bin/cyborgnerd

tw:
	tailwindcss -i static/input.css -o static/output.css --watch
dev:
	go run .
css:
	tailwindcss -i static/input.css -o static/output.css --minify
binary:
	CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o "$(BIN)" .
build: css binary
test:
	go test ./...
