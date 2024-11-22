BINARY_NAME=api

run:
	go run main.go -listenaddr=:8080

test:
	go test -v ./...

build:
	$(GO) build -o $(BINARY_NAME)

clean:
	rm -f $(BINARY_NAME)