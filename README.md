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

# Logic
- The Alexa skill lambda will send request to lambda that processes chatGPT prompts via SQS, the Alexa skill lambda will then poll the response SQS and return the response if a message is avaliable in <7 seconds.

> Due to the Alexa skill response constraint of 8 seconds following logic has been applied
- if the Alexa skill doesn't receive a SQS message when polling the response SQS within ~7 seconds, it will return 'your response will be available momentarily', to not end the Alexa skill session.
- Querying the Alexa skill with 'Last Response', the Alexa Skill lambda will immediately poll the response SQS to retrieve the delayed response and output the prompt with the timestamp of response time

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
  export HANDLER=main
```

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

### Requirements
- <b>OPENAI API KEY</b>
  * please set environment variables for your OPENAI API key
    > export API_KEY=123456


- <b>Create a S3 Bucket on your AWS Account</b>
  * Set envrionment variable of the S3 Bucket name you have created [this is where AWS SAM]
    > export S3_BUCKET_NAME=bucket_name


### Deployment Steps
1. Create a new Alexa skill with a name of your choice


2. Set the Alexa skill invocation with a phrase i.e. 'My question'


3. Set built-in invent invocations to their relevant phrases i.e. 'help', 'stop', 'cancel', etc.


4. Create a new Intent named 'AutoCompleteIntent'


5. Add a new Alexa slot to this Intent and name it 'prompt' with type AMAZON.SearchQuery'


6. Add invocation phrase for the 'AutoCompleteIntent' with value 'question {prompt}'


7. Deploy the stack to your AWS account.
> sam build  && sam deploy --stack-name chat-gpt --s3-bucket $S3_BUCKET_NAME --parameter-overrides 'ApiKey=$API_KEY' --capabilities CAPABILITY_IAM


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
