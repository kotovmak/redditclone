.PHONY: build
build:
	go build -v -o ./build/apiserver/apiserver ./cmd/apiserver/main.go

.PHONY: test
test:
	go test -v -race -timeout 30s ./...

.PHONY: run
run: build
	./build/apiserver/apiserver

.DEFAULT_GOAL := run
