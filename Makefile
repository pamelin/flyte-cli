build:
	dep ensure
	go test ./...
	go build -o flyte

install:
	dep ensure
	go test ./...
	go build -o $(GOPATH)/bin/flyte
