package main

import (
	"context"
	"os"
	"strconv"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/jackmcguire1/alexa-chatgpt/internal/api"
	"github.com/jackmcguire1/alexa-chatgpt/internal/dom/chatmodels"
	pkginit "github.com/jackmcguire1/alexa-chatgpt/internal/pkg/init"
	"github.com/jackmcguire1/alexa-chatgpt/internal/pkg/queue"
	"go.opentelemetry.io/contrib/instrumentation/github.com/aws/aws-lambda-go/otellambda"
	"go.opentelemetry.io/contrib/instrumentation/github.com/aws/aws-lambda-go/otellambda/xrayconfig"
)

func main() {
	logger := pkginit.SetupLogger()
	ctx := context.Background()
	tracer := pkginit.SetupTracing(ctx, logger)
	defer tracer.Shutdown(ctx)

	resources := pkginit.InitializeResources()
	svc := chatmodels.NewClient(resources)
	pollDelay, _ := strconv.Atoi(os.Getenv("POLL_DELAY"))

	h := api.Handler{
		ChatGptService:  svc,
		RequestsQueue:   queue.NewQueue(os.Getenv("REQUESTS_QUEUE_URI")),
		ResponsesQueue:  queue.NewQueue(os.Getenv("RESPONSES_QUEUE_URI")),
		PollDelay:       pollDelay,
		Logger:          logger,
		Model:           pkginit.GetDefaultChatModel(resources),
		ImageModel:      pkginit.GetDefaultImageModel(resources),
		RandomNumberSvc: api.NewRandomNumberGame(100),
		BattleShips:     api.NewBattleShipSetup(),
		AnimalGame:      api.NewAnimalGame(),
	}
	lambda.Start(otellambda.InstrumentHandler(h.Invoke, xrayconfig.WithRecommendedOptions(tracer)...))
}
