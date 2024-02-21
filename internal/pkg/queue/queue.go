package queue

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/jackmcguire1/alexa-chatgpt/internal/pkg/utils"
)

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

	data := utils.ToJSON(i)

	_, err := q.client.SendMessage(ctx, &sqs.SendMessageInput{
		MessageBody:  &data,
		QueueUrl:     &q.queueUri,
		DelaySeconds: 0,
	})

	return err
}

func (q *Queue) PullMessage(ctx context.Context, wait int) ([]byte, error) {
	resp, err := q.client.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
		QueueUrl:            &q.queueUri,
		MaxNumberOfMessages: 1,
		WaitTimeSeconds:     int32(wait),
	})
	if err != nil {
		return nil, err
	}

	if len(resp.Messages) == 0 {
		return nil, nil
	}

	_, err = q.client.DeleteMessage(ctx, &sqs.DeleteMessageInput{
		QueueUrl:      &q.queueUri,
		ReceiptHandle: resp.Messages[0].ReceiptHandle,
	})
	if err != nil {
		return nil, err
	}

	return []byte(*resp.Messages[0].Body), nil
}

func (q *Queue) Purge(ctx context.Context) error {
	_, err := q.client.PurgeQueue(ctx, &sqs.PurgeQueueInput{QueueUrl: &q.queueUri})
	return err
}
