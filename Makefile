SHELL=/usr/bin/env bash -o pipefail

all: lint test build
.PHONY: all

build:
	go build -o ./out/tist ./cmd/...
.PHONY: build

install:
	go install ./cmd/...
.PHONY: install

clean:
	rm -rf ./tist
.PHONY: clean

test:
	go test -race ./...
.PHONY: test

lint:
	gofmt -s -l $(shell go list -f '{{ .Dir }}' ./... ) | grep ".*\.go"; if [ "$$?" = "0" ]; then gofmt -s -d $(shell go list -f '{{ .Dir }}' ./... ); exit 1; fi
	go vet ./...
.PHONY: lint

format:
	gofmt -s -w $(shell go list -f '{{ .Dir }}' ./... )
.PHONY: format
