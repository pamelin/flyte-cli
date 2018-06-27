build:
	dep ensure -v
	go test ./...
	go build -o flyte

install:
	dep ensure -v
	go test ./...
	go build -o $(GOPATH)/bin/flyte

justdoit:
	go build -o $(GOPATH)/bin/flyte
