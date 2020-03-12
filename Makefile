.PHONY: build clean deploy

build:
	env GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o bin/hello lambda/hello/hello.go
	env GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o bin/messages lambda/messages/messages.go
	env GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o bin/redirectUrl lambda/redirectUrl/redirectUrl.go
	env GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o bin/notifications lambda/notifications/notifications.go
	env GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o bin/leaderboard lambda/leaderboard/leaderboard.go
	env GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o bin/interactive lambda/interactive/main.go

clean:
	rm -rf ./bin

deploy: clean build
	sls deploy --verbose
