.PHONY: build clean deploy

build:
	env GOOS=linux go build -ldflags="-s -w" -o bin/hello hello/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/dynamodb-example dynamodb-example/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/incoming-webhook incoming-webhook/main.go

clean:
	rm -rf ./bin

deploy: clean build
	sls deploy --verbose
