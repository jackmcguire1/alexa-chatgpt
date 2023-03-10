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

# Examples

<p align="center">
  <img src="./images/image.png" width="350" height="500" title="Random Phrase">
  <img src="./images/image_1.png" width="350" height="500" title="Alexa is inferior">
</p>

## SETUP
> How to configure your Alexa Skill

> Please Configure the Makefile with your own available S3 bucket

1. Create a new Alexa skill with a name of your choice

2. Set Alexa skill invocation with a sentence of your choice or i.e. 'sentence completer'

3. Set built-in invent invocations to their relevant phrases i.e. 'help', 'stop', 'cancel', etc.

4. Set the fallback intent invocation phrase to be 'random phrase'

5. Create new Intent named 'AutoCompleteIntent'

6. Add invocation phrase for the 'AutoCompleteIntent' with value 'complete the sentence {prompt}'

7. deploy chatgpt lambda and take note of the ARN of the lambda

8. Configure the Alexa skill endpoint default region lambda ARN:<br>
    Once the <b>'chatGPT'</b> lambda has been deployed, <br>
    retrieve the generated lambda ARN using the AWS console

9. Begin testing your Alexa skill by querying for 'sentence completer' or your chosen invocation phrase

10. Query Alexa 'complete the sentence {your sentence here}'

11. Query Alexa 'random phrase'!

12. Tell Alexa to 'stop'

13. <b>Testing complete!</b>

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