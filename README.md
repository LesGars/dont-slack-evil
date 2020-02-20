# Don't Slack Evil

Don't Slack Evil is the submission of the team "Les Gars" for the [2020 Slack App virtual hackathon](https://slackapponlinehackathon.splashthat.com/)

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