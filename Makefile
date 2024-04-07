build:
	@go build -o bin/gowallet

run: build
	@bin/gowallet

test:
	@go test -v ./...