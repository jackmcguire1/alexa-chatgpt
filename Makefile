build-ChatGPTFunction:
	GOOS=linux GOARCH=amd64 go build -o main github.com/jackmcguire1/alexa-chatgpt/cmd/alexa/
	cp ./main $(ARTIFACTS_DIR)/
  
build-ChatGPTRequests:
	GOOS=linux GOARCH=amd64 go build -o main github.com/jackmcguire1/alexa-chatgpt/cmd/sqs/
	cp ./main $(ARTIFACTS_DIR)/
  
  