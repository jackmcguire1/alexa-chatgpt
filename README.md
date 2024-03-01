# Alexa-ChatGPT

> This repository contains the Alexa skill serverless backend to prompt generative ai LLM models

[git]: https://git-scm.com/
[golang]: https://golang.org/
[modules]: https://github.com/golang/go/wiki/Modules
[golint]: https://github.com/golangci/golangci-lint
[aws-cli]: https://docs.aws.amazon.com/cli/latest/userguide/cli-chap-install.html
[aws-cli-config]: https://docs.aws.amazon.com/cli/latest/userguide/cli-chap-configure.html
[aws-sam-cli]: https://github.com/awslabs/aws-sam-cli

[![Go Report Card](https://goreportcard.com/badge/github.com/jackmcguire1/alexa-chatgpt)](https://goreportcard.com/report/github.com/jackmcguire1/alexa-chatgpt)
[![codecov](https://codecov.io/gh/jackmcguire1/alexa-chatgpt/branch/main/graph/badge.svg)](https://codecov.io/gh/jackmcguire1/alexa-chatgpt)

# Logic

- A user prompts the Alexa skill.
- The Alexa skill will invoke the assigned Lambda with an 'AutoComplete' Intent.
- The Lambda will push the user prompt to a SQS.
- The request lambda will be invoked with the SQS message and begin to process the user's prompt via the chosen chat model [OpenAI ChatGPT , Google Gemini] and put the response onto a seperate SQS.
- Meanwhile the Alexa skill lambda will be polling the response SQS in order to return the response for the prompt.

> [!CAUTION]
> Due to the Alexa skill idle lambda response constraint of ~8 seconds, the following logic has been applied.

- If the Alexa skill does not poll a message from the queue within ~7 seconds, users will be given a direct response of 'your response will be available shortly!', this is too avoid the Alexa skill session from expiring.

- Querying the Alexa skill with 'Last Response', the lambda will immediately poll the response SQS to retrieve the delayed response and output the prompt with the timestamp of response time

## Supported Models

> [!NOTE]
> Users are able to change which chat model is in use

### OpenAI ChatGPT

- user's can select this by prompting 'use gpt'

### Google's GenerativeAI Gemini

- user's can select this by prompting 'use gemini'

## Alexa Intents

> The Alexa Intents or phrases to interact with the Alexa Skill

- AutoComplete

  > the intent used to prompt the LLM models

- Model

  > Allows users to select LLM model to use

- Last Response

  > Fetch delayed LLM response to user's prompt

- Cancel

  > Force Alexa to await for next intent

- Stop

  > Terminate Alexa skill session

- Help
  > List all avalible interactions or intents

# Infrastructure

  <img src="./images/infra.png">

# Examples

<p align="center">
  <img src="./images/image.png" width="400" height="500" title="Example question">
  </p>

## SETUP

> How to configure your Alexa Skill

### Environment

> we use handler env var to name the go binary either 'main' or 'bootstrap' for AL2.Provided purposes, devs should use 'main'

```shell
  HANDLER=main
  OPENAI_API_KEY=xxx
  GEMINI_API_KEY={base64 service account json}
```

### Prerequisites

- [Git][git]
- [Go 1.21][golang]+
- [golangCI-Lint][golint]
- [AWS CLI][aws-cli]
- [AWS SAM CLI][aws-sam-cli]

### [AWS CLI Configuration][aws-cli-config]

> Make sure you configure the AWS CLI

- AWS Access Key ID
- AWS Secret Access Key
- Default region 'us-east-1'

```shell
aws configure
```

### Requirements

- <b>OPENAI API KEY</b>

  - please set environment variables for your OPENAI API key
    > export API_KEY=123456

- <b>Create a S3 Bucket on your AWS Account</b>
  - Set envrionment variable of the S3 Bucket name you have created [this is where AWS SAM]
    > export S3_BUCKET_NAME=bucket_name

### Deployment Steps

1. Create a new Alexa skill with a name of your choice

2. Set the Alexa skill invocation with a phrase i.e. 'My question'

3. Set built-in invent invocations to their relevant phrases i.e. 'help', 'stop', 'cancel', etc.

4. Create a new Intent named 'AutoCompleteIntent'

5. Add a new Alexa slot to this Intent and name it 'prompt' with type AMAZON.SearchQuery'

6. Add invocation phrase for the 'AutoCompleteIntent' with value 'question {prompt}'

7. Deploy the stack to your AWS account.

```
  export ARCH=GOARCH=arm64
  export LAMBDA_RUNTIME=provided.al2023
  export LAMBDA_HANDLER=bootstrap
  export LAMBDA_ARCH=arm64
```

```
sam build --parameter-overrides Runtime=$LAMBDA_RUNTIME Handler=$LAMBDA_HANDLER Architecture=$LAMBDA_ARCH
```

```
sam deploy --stack-name chat-gpt --s3-bucket $S3_BUCKET --parameter-overrides Runtime=$LAMBDA_RUNTIME Handler=$LAMBDA_HANDLER Architecture=$LAMBDA_ARCH OpenAIApiKey=$OPENAI_API_KEY GeminiApiKey=$GEMINI_API_KEY --capabilities CAPABILITY_IAM

```

9. Once the stack has deployed, make note of lambda ARN from the 'ChatGPTLambdaArn' field, from the the output of

   > sam list stack-outputs --stack-name chat-gpt

10. Apply this lambda ARN to your 'Default Endpoint' configuration within your Alexa skill, i.e. 'arn:aws:lambda:us-east-1:123456789:function:chatGPT'

11. Begin testing your Alexa skill by querying for 'My question' or your chosen invocation phrase, Alexa should respond with "Hi, let's begin our conversation!"

12. Query Alexa 'question {your sentence here}'

    > Note the OpenAI API may take longer than 8 seconds to respond, in this scenario Alexa will tell you your answer will be ready momentarily, simply then ask Alexa 'last response'

13. Tell Alexa to 'stop'

14. <b>Testing complete!</b>

## Contributors

This project exists thanks to **all** the people who contribute.

## Donations

All donations are appreciated!

[![Donate](https://img.shields.io/badge/Donate-PayPal-green.svg)](http://paypal.me/crazyjack12)
