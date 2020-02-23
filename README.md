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
3) Deploy: `make deploy`

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