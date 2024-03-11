package bucket

import (
	"bytes"
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type FilePersistance interface {
	Put(ctx context.Context, reqId string, fileName string, prefix string, data []byte) (string, error)
	FolderName() string
}

type Bucket struct {
	Name string
}

func (b *Bucket) Put(ctx context.Context, reqId string, fileName string, prefix string, data []byte) (string, error) {
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("us-east-1"))
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	svc := s3.NewFromConfig(cfg)
	_, err = svc.PutObject(ctx, &s3.PutObjectInput{
		Bucket: &b.Name,
		Key:    aws.String(prefix + reqId + "/" + fileName),
		Body:   bytes.NewReader(data),
	})
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("https://s3.amazonaws.com/%s/%s%s/%s",
		b.Name,
		prefix,
		reqId,
		fileName,
	), nil
}

func (b *Bucket) FolderName() string {
	return b.Name
}
