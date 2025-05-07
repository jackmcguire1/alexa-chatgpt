package queue

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/jackmcguire1/alexa-chatgpt/internal/pkg/utils"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

var tracer = otel.Tracer("sqs")

var EmptyMessageErr = fmt.Errorf("no messages found")

type PullPoll interface {
	Purge(context.Context) error
	PullMessage(context.Context, int) ([]byte, error)
	PushMessage(context.Context, any) error
}

type Queue struct {
	client   *sqs.Client
	queueUri string
}

func NewQueue(queueUri string) *Queue {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"))
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	return &Queue{
		client:   sqs.NewFromConfig(cfg),
		queueUri: queueUri,
	}
}

func (q *Queue) PushMessage(ctx context.Context, i any) error {
	ctx, span := tracer.Start(ctx, "SendMessage")
	defer span.End()
	span.SetAttributes(attribute.String("queue.uri", q.queueUri))

	data := utils.ToJSON(i)

	_, err := q.client.SendMessage(ctx, &sqs.SendMessageInput{
		MessageBody:  &data,
		QueueUrl:     &q.queueUri,
		DelaySeconds: 0,
	})
	if err != nil {
		span.RecordError(err)
		return err
	}

	return nil
}

func (q *Queue) PullMessage(ctx context.Context, wait int) ([]byte, error) {
	ctx, span := tracer.Start(ctx, "PullMessage")
	defer span.End()
	span.SetAttributes(attribute.String("queue.uri", q.queueUri))

	resp, err := q.client.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
		QueueUrl:            &q.queueUri,
		MaxNumberOfMessages: 1,
		WaitTimeSeconds:     int32(wait),
	})
	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	if len(resp.Messages) == 0 {
		return nil, nil
	}

	// Start a child span for message deletion
	deleteCtx, deleteSpan := tracer.Start(ctx, "DeleteMessage")
	defer deleteSpan.End()
	span.SetAttributes(
		attribute.String("queue.uri", q.queueUri),
		attribute.String("message-receipt", *resp.Messages[0].ReceiptHandle),
	)

	_, err = q.client.DeleteMessage(deleteCtx, &sqs.DeleteMessageInput{
		QueueUrl:      &q.queueUri,
		ReceiptHandle: resp.Messages[0].ReceiptHandle,
	})
	if err != nil {
		deleteSpan.RecordError(err)
		return nil, err
	}

	return []byte(*resp.Messages[0].Body), nil
}

func (q *Queue) Purge(ctx context.Context) error {
	ctx, span := tracer.Start(ctx, "Purge")
	defer span.End()
	span.SetAttributes(attribute.String("queue.uri", q.queueUri))

	_, err := q.client.PurgeQueue(ctx, &sqs.PurgeQueueInput{QueueUrl: &q.queueUri})
	return err
}
