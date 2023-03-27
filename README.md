# Alexa-ChatGPT
> This repository contains the Alexa skill to use the OpenAI API

[git]:    https://git-scm.com/
[golang]: https://golang.org/
[modules]: https://github.com/golang/go/wiki/Modules
[golint]: https://github.com/golangci/golangci-lint
[aws-cli]: https://docs.aws.amazon.com/cli/latest/userguide/cli-chap-install.html
[aws-cli-config]: https://docs.aws.amazon.com/cli/latest/userguide/cli-chap-configure.html
[aws-sam-cli]: https://github.com/awslabs/aws-sam-cli


[![Go Report Card](https://goreportcard.com/badge/github.com/jackmcguire1/alexa-chatgpt)](https://goreportcard.com/report/github.com/jackmcguire1/alexa-chatgpt)
[![codecov](https://codecov.io/gh/jackmcguire1/alexa-chatgpt/branch/main/graph/badge.svg)](https://codecov.io/gh/jackmcguire1/alexa-chatgpt)


# Infrastructure
  <img src="./images/infra.png">

# Examples

<p align="center">
  <img src="./images/image.png" width="400" height="500" title="Example question">
  </p>

## SETUP
> How to configure your Alexa Skill

> Please Configure the Makefile with your own available S3 bucket

1. Create a new Alexa skill with a name of your choice


2. Set the Alexa skill invocation with a phrase i.e. 'My question'


3. Set built-in invent invocations to their relevant phrases i.e. 'help', 'stop', 'cancel', etc.


4. Create a new Intent named 'AutoCompleteIntent'


5. Add a new Alexa slot to this Intent and name it 'prompt' with type AMAZON.SearchQuery'


6. Add invocation phrase for the 'AutoCompleteIntent' with value 'question {prompt}'


7. Deploy the stack to your AWS account.
> sam build  && sam deploy --stack-name chat-gpt --s3-bucket $S3_BUCKET_NAME --parameter-overrides 'ApiKey=$API_KEY' --capabilities CAPABILITY_IAM


9. Once the stack has deployed, make note of lambda ARN from the 'ChatGPTLambdaArn' field, from the the output of > sam list stack-outputs --stack-name chat-gpt


10. Apply this lambda ARN to your 'Default Endpoint' configuration within your Alexa skill, i.e. 'arn:aws:lambda:us-east-1:123456789:function:chatGPT'


11. Begin testing your Alexa skill by querying for 'My question' or your chosen invocation phrase, Alexa should respond with "Hi, let's begin our conversation!"


12. Query Alexa 'question {your sentence here}'


13. > Note the OpenAI API may take longer than 8 seconds to respond, in this scenario Alexa will tell you your answer will be ready momentarily, simply then ask Alexa 'last response'


14. Tell Alexa to 'stop'


15. <b>Testing complete!</b>

## Development

To develop `how-old-is` or interact with its source code in any meaningful way, be
sure you have the following installed:

### Prerequisites

- [Git][git]
- [Go 1.20][golang]+
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

## Contributors

This project exists thanks to **all** the people who contribute.

## Donations
All donations are appreciated!

[![Donate](https://img.shields.io/badge/Donate-PayPal-green.svg)](http://paypal.me/crazyjack12)
