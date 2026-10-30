.PHONY: run test check build
run:
	go run ./cmd/server
test:
	go test -race ./...
check:
	test -z "$$(gofmt -l cmd internal)"
	go vet ./...
	go test -race ./...
build:
	go build -trimpath -o bin/server ./cmd/server
