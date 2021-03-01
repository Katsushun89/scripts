
NAME := parallel_exec_script

test:
	go test

build:
	for arch in amd64 arm64; do \
		CGO_ENABLED=0 GOOS=linux GOARCH=$$arch go build -o bin/$$arch/$(NAME); \
	done

clean: 
	rm bin/* -rf

.PHONY: test build clean

