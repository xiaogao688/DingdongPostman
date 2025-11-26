.PHONY: build run test clean

build:
	go build -o bin/dingdong-postman main.go

run:
	go run main.go

test:
	go test -v ./...

clean:
	rm -rf bin/
	go clean

fmt:
	go fmt ./...

lint:
	golangci-lint run

vet:
	go vet ./...

