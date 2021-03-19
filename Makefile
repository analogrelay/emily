export GO111MODULE=on

.PHONY: all
all: emcee

emcee:
	go build -o bin/emcee cmd/emcee

test:
	go test ./...