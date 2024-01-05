build:
	@go build -o bin/cx

run: build
	@./bin/cx

test:
	@go test -v ./...
