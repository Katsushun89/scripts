
NAME := scripts

test:
	go test -v ./...

build:
	for arch in amd64 arm64; do \
		CGO_ENABLED=0 GOOS=linux GOARCH=$$arch go build -o bin/$$arch/$(NAME); \
	done

clean: 
	rm bin/* -rf

.PHONY: test build clean

