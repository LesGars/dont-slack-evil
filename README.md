[![Go Report Card](https://goreportcard.com/badge/github.com/gjgd/dont-slack-evil)](https://goreportcard.com/report/github.com/gjgd/dont-slack-evil)
![Continuous production](https://github.com/gjgd/dont-slack-evil/workflows/Continuous%20production/badge.svg)

# Don't Slack Evil

Don't Slack Evil is the submission of the team "Les Gars" for the [2020 Slack App virtual hackathon](https://slackapponlinehackathon.splashthat.com/)

# Getting started

1) Setup the [serverless framework](https://github.com/serverless/serverless)
2) Set up shared AWS credentials:
  - Ask Cyril for the shared AWS credentials
  - Add to your AWS credentials file (`~/.aws/credentials`) the following:
  ```
  [dont-slack-evil-hackaton]
  aws_access_key_id=SECRET_FROM_CYRIL
  aws_secret_access_key=SECRET_FROM_CYRIL
  ```
  - Add to your aws config file (`~/.aws/config`) the following:
  ```
  [profile dont-slack-evil-hackaton]
  region=us-east-1
  ```
3) Set up secret file: `cp example.secrets.dev.yml secrets.dev.yml` and fill the secrets
4) Deploy: `make deploy`

# Requirements
* [golang](https://golang.org/dl/) >= 1.13
* [direnv](https://direnv.net/) (if you are not using VSCode)

# Updating env variables

Currently if you need to update environment variables you need to
- update `example.secrets.dev.yml`
- update `secrets.dev.yml`, the actual env file
- update `serverless.yml` to inject the env variable where it is needed
- update the "Create secrets file" section of `continuous_production.yml` to create the env file in CI

# Setting up your Go environment

## If you are using VSCode

* Install the Go extension. When prompted, install the necessary *Go Tools* (`golint`, `goreturns`, `gopls`)
* Add these properties to your `settings.json` file:
```json
{
    "go.useLanguageServer": true,
    "go.toolsEnvVars": {
            "GO111MODULES": "on",
            "GOFLAGS": "-mod=vendor"
    }
}
```
The first property, `go.useLanguageServer`, will activate the Language Server for Go, which will improve features such as autocompletion, when using modules.

The second property, `go.toolsEnvVars`, will export 2 env vars when code analysis is performed by VSCode:

1. `GO111MODULES=on`: this will activate the use of [Golang modules](https://blog.golang.org/using-go-modules)
2. `GOFLAGS=-mod=vendor`: references to external packages will be looked up in the `/vendor` gitignored folder

You can make use of the following extra properties:

* `go.testOnSave: true` for triggering `go test` on each save (if autosave is disabled) or CTRL+S (if autosave is enabled)
* `go.testOnSave: false` for deactivating build on each save

### Setup tests

Under Debug > open configurations, add the following ENV variables
```
"env": {
  "DYNAMODB_TABLE": "dont-slack-evil-test-messages",
  "AWS_REGION": "us-east-1"
},
```
Under Code > Preferences > Settings, search for "testEnvVars" and add your vars again under settings
```
"settings": {
  "go.testEnvVars": {
    "DYNAMODB_TABLE": "dont-slack-evil-test-messages",
    "AWS_REGION": "us-east-1"
  },
}
```

## If you are using a regular text editor
* `direnv` will execute the `.envrc` file and leverage the 2 env variables
* Install `golint` (the linter) and `goreturns` (the formatter)
```
go get -u golang.org/x/lint/golint
go get -u github.com/sqs/goreturns
```

* Before committing your code, run the formatter, and then, run the linter and correct any error it raises
```
go fmt ./..
$HOME/go/bin/golint ./..
```
**WARNING:** The linter will also show errors from the `/vendor` directory

## Adding/removing a dependency

* If you are using VSCode, the added dependency will be downloaded in the background. Otherwise, you will have to run
```
go get -u <Dependency_Path>
```
* Delete the `vendor/` folder, cd into the project's root dir and run these 2 commands:
```
# Recreate the vendor/ folder
go mod vendor

# Update go.mod and go.sum
go mod tidy
```

## Running tests

Project-level tests can be run with this single command:

```
go test ./...
```
You can append the `-v` flag for verbose input

# Useful links

## Serverless AWS Lambda in Golang

https://github.com/serverless/examples

## DynamoDB

https://github.com/awsdocs/aws-doc-sdk-examples/tree/456068072ee2ee696b2aac4724c925171c2bb2ff/go/example_code/dynamodb

# Useful commands

```bash
# Deploy
serverless deploy function -f webhook
# OR
serverless deploy
```
```bash
# Logs the activity of a lambda
serverless logs -t -f hello
```
```bash
# Encrypt / Decrypt secrets
serverless encrypt --stage dev --password 'your-password'
serverless decrypt --stage dev --password 'your-password'
```
```bash
# Remove all serverless services
serverless remove
```

If `make deploy` fails:

```golang
go get -u github.com/aws/aws-lambda-go/lambda
```
