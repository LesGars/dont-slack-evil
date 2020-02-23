[![Go Report Card](https://goreportcard.com/badge/github.com/gjgd/dont-slack-evil)](https://goreportcard.com/report/github.com/gjgd/dont-slack-evil)
# Don't Slack Evil

Don't Slack Evil is the submission of the team "Les Gars" for the [2020 Slack App virtual hackathon](https://slackapponlinehackathon.splashthat.com/)

# Requirements
* [golang](https://golang.org/dl/) >= 1.13
* [direnv](https://direnv.net/) (if you are not using VSCode)

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

## Adding a new dependency

* If you are using VSCode, the dependency will be downloaded in the background. Else, you will have to run
```
go get -u <Dependency_Path>
```
* Delete the `vendor/` folder and run
```
go mod vendor
```
at the root of the project

# Getting started

1) Setup the [serverless framework](https://github.com/serverless/serverless)
2) Deploy: `make deploy`

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

