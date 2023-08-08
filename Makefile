build-ChatGPTRequests:
	GOOS=linux ${ARCH} go build -o ${HANDLER} github.com/jackmcguire1/alexa-chatgpt/cmd/sqs/
	cp ./${HANDLER} $(ARTIFACTS_DIR)/.
  
  
