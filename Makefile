build:
	go build -o ./bin/cx

run:
	./bin/cx

test:
	go test -v ./...
