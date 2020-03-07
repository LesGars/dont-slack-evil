.PHONY: build clean deploy

build:
	env GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o bin/hello lambda/hello/hello.go
	env GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o bin/dynamodbexample lambda/dynamodbexample/dynamodbexample.go
	env GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o bin/messages lambda/messages/messages.go

clean:
	rm -rf ./bin

deploy: clean build
	sls deploy --verbose
