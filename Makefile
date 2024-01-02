build:
	@go build -o bin/cx

run:
	@go build -o bin/cx; ./bin/cx

test:
	@go test -v ./...
