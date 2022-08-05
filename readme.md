# AWS SNS Go Lambda Template


- [AWS SNS Lambda](aws-sns-go-lambda-template)
  - [Technologies](#technologies)
  - [Overview](#overview)
  - [Local Development](#local-development)
    - [Installation Pre-requisites](#installation-pre-requisites)
    - [Installing SAM CLI](#installing-sam-clie)
    - [Configuring AWS Credentials](#configuring-aws-credentials)
    - [Debugging](#debugging)
        - [Configuring Delve](#configure-delve)
        - [Direct Launch](#direct-launch)
  - [Deployment](#deployment)

## Technologies
* Go 1.x (the minor version is specified by AWS)
* AWS Lambda
* AWS CloudFormation

## Overview

This is a sample of a working AWS Lambda written in Go that uses the Go AWS-SDK to interact with AWS SNS topics. The Lambda is intended to show how to trigger a lambda from an inbound subscription to an SNS Topic, as well as publish new messages to an outbound SNS Topic.

`SnsDebugger` simply outputs the incoming data to the console and once deployed this output can be found in the CloudWatch logs. This will help to debug and examine the incoming data.

`SnsProcessor` reads the incoming sns messages and uses `gjson` to extract data from the incoming SNS message. It then creates and publishes a new SNS message to an outbound topic. The incoming message is included as the data of the outbound message, along with custom MessageAttributes.

## Local Development

### Installation Pre-requisites
* GoLang 1.x
* AWS CLI - https://aws.amazon.com/cli/
* AWS SAM CLI - https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/serverless-sam-cli-install.html
* VSCode (optional)
* Docker

### Installing SAM CLI

In order to run and test your lambda locally, you will need to install the SAM CLI. AWS provides guidance here: https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/serverless-sam-cli-install.html

### Configuring AWS Credentials

For the CLI to run you'll need to configure it with AWS credentials to run under. See the documentation here: https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/serverless-getting-started-set-up-credentials.html

The recommendation is to skip to the `Not using the AWS CLI` section and create a credentials file that includes the AWS access key and secret.

### Invoking the Lambda locally
Invoke the lambda function with the following:
```
sam build
sam local invoke "SnsProcessor" --event events/sample-sns-event.json --env-vars .\env\env-dev.json
```


### Debugging


#### Configuring Delve
In order to debug locally, you'll want to make sure delve is installed. From the `./handlers` folder run the following:
```
mkdir delve
GOARCH=amd64 GOOS=linux go install github.com/go-delve/delve/cmd/dlv@late
GOARCH=amd64 GOOS=linux go build -o delve/dlv github.com/go-delve/delve/cmd/dlv
GOARCH=amd64 GOOS=linux go build -o dlv github.com/go-delve/delve/cmd/dlv
```
In the next section we'll pass the location of delve into SAM invoke command so we can locally debug the code

#### Direct Launch 
The included `launch.json` file in the `.vscode` folder allows you to invoke the method locally via AWS SAM CLI in debug mode and attach the debugger. Environment variables are loaded from a file specified in the command line.

First invoke the lambda via the AWS SAM CLI terminal and include the `-d` debug option and port number. Also specify the template and environemnt variable file to use.

`>sam local invoke "SnsDebugger" --event ./events/sample-sns-event.json --env-vars=./env/env-dev.json -d 9999 --debugger-path handlers/delve/ --debug-args "-delveAPI=2"`

The port number should match the one defined in the `launch.json` config:

```json
{
    "configurations": [
        {
            "name": "Debug SAM Lambda",
            "type": "go",
            "request": "attach",
            "mode": "remote",
            "remotePath": "",
            "port": 9999,

```

Once the lambda is running the terminal will report the debugger is listening:

`API server listening at: [::]:9999`

From there select Run & Debug from the left menu and select one of the `Debug SAM Lambda` profile.

## Deployment 

To deploy this project to AWS run the `sam deploy` command and follow the guided prompts.