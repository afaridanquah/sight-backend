BINARY=verifylab-service

hello:
	echo "Hello"

tidy:
	go mod tidy
	go mod vendor

build:
	@go build -o bin/${BINARY}

run: build
	@./bin/api/${BINARY}

clean: 
	go clean
