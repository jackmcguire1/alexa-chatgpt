build-ChatGPTFunction:
	CGO_ENABLED=0 GOOS=linux ${ARCH} go build -o ${HANDLER} github.com/jackmcguire1/alexa-chatgpt/cmd/alexa/
	cp ./${HANDLER} $(ARTIFACTS_DIR)/.
  
build-ChatGPTRequests:
	CGO_ENABLED=0 GOOS=linux ${ARCH} go build -o ${HANDLER} github.com/jackmcguire1/alexa-chatgpt/cmd/sqs/
	cp ./${HANDLER} $(ARTIFACTS_DIR)/.
  
  
