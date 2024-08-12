BINARY=verifylab-backend

hello:
	echo "Hello"

build:
	@go build -o bin/${BINARY}

run: build
	@./bin/api/${BINARY}

clean: 
	go clean
