# Alexa-ChatGPT
> This repository contains the Alexa skill to use the OpenAI API

[git]:    https://git-scm.com/
[golang]: https://golang.org/
[modules]: https://github.com/golang/go/wiki/Modules
[golint]: https://github.com/golangci/golangci-lint
[aws-cli]: https://docs.aws.amazon.com/cli/latest/userguide/cli-chap-install.html
[aws-cli-config]: https://docs.aws.amazon.com/cli/latest/userguide/cli-chap-configure.html
[aws-sam-cli]: https://github.com/awslabs/aws-sam-cli

## SETUP
> How to configure your Alexa Skill

> Please Configure the Makefile with your own available S3 bucket

1. Create a new Alexa skill with a name of your choice

2. Set Alexa skill invocation method as 'chatGPT prompt'

3. Set built-in invent invocations to their relevant phrases i.e. 'help', 'stop', 'cancel', etc.

4. Set a random phrase for the built-in fallback intent, i.e. 'bumbaclart'

5. Create new Intent named 'ChatGPTIntent'

6. Create new intent slot named 'name', with <b>SLOT TYPE</b> 'AMAZON.FirstName'

7. Add invocation phrase for the 'ChatGPTIntent' intent with phrase 'chatgpt prompt'

8. Configure slot values for the 'AMAZON.FirstName' <b>SLOT TYPE</b> i.e. 'Bob'

9. Edit the 'internal/dom/age.go' file to feature your name and DoB in RFC3339 format

10. Package and deploy how-old-is lambda

11. Configure Alexa skill endpoint lambda ARN:<br>
    Once the <b>'how-old-is'</b> lambda has been deployed, <br>
    retrieve the generated lambda ARN using the AWS console or<br>
    one of the describe stack methods found above.<br>
    input the lambda <b>ARN</b> as the default endpoint of your Alexa skill,<br>
    within your Alexa development console!

12. Begin testing your Alexa skill by querying for 'chatgpt prompt'

13. Query Alexa 'chatgpt prompt {question}'

14. Query Alexa 'bumbaclart' or your fallback invocation phrase!

15. Tell Alexa to 'stop'

16. <b>Testing complete!</b>

## Development

To develop `how-old-is` or interact with its source code in any meaningful way, be
sure you have the following installed:

### Prerequisites

- [Git][git]
- [Go 1.20][golang]+
- [golangCI-Lint][golint]
- [AWS CLI][aws-cli]
- [AWS SAM CLI][aws-sam-cli]

>You will need to activate [Modules][modules] for your version of [GO][golang],

> by setting the `GO111MODULE=on` environment variable set

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