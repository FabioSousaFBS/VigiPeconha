package storage

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Storage interface {
	Upload(
		ctx context.Context,
		key string,
		body io.Reader,
		contentType string,
	) error

	Delete(
		ctx context.Context,
		key string,
	) error

	PresignGetURL(
		ctx context.Context,
		key string,
		expires time.Duration,
	) (string, error)
}

type R2Storage struct {
	client *s3.Client
	bucket string
}

func NewR2Storage(
	ctx context.Context,
	accountID string,
	accessKeyID string,
	secretAccessKey string,
	bucket string,
) (*R2Storage, error) {
	awsConfig, err := config.LoadDefaultConfig(
		ctx,
		config.WithRegion("auto"),
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(
				accessKeyID,
				secretAccessKey,
				"",
			),
		),
	)

	if err != nil {
		return nil, fmt.Errorf(
			"erro ao carregar configuração R2: %w",
			err,
		)
	}

	endpoint := fmt.Sprintf(
		"https://%s.r2.cloudflarestorage.com",
		accountID,
	)

	client := s3.NewFromConfig(
		awsConfig,
		func(options *s3.Options) {
			options.BaseEndpoint = aws.String(endpoint)
			options.UsePathStyle = true
		},
	)

	return &R2Storage{
		client: client,
		bucket: bucket,
	}, nil
}

func (s *R2Storage) Upload(
	ctx context.Context,
	key string,
	body io.Reader,
	contentType string,
) error {
	_, err := s.client.PutObject(
		ctx,
		&s3.PutObjectInput{
			Bucket:      aws.String(s.bucket),
			Key:         aws.String(key),
			Body:        body,
			ContentType: aws.String(contentType),
		},
	)

	if err != nil {
		return fmt.Errorf(
			"erro ao enviar arquivo para R2: %w",
			err,
		)
	}

	return nil
}

func (s *R2Storage) Delete(
	ctx context.Context,
	key string,
) error {
	_, err := s.client.DeleteObject(
		ctx,
		&s3.DeleteObjectInput{
			Bucket: aws.String(s.bucket),
			Key:    aws.String(key),
		},
	)

	if err != nil {
		return fmt.Errorf(
			"erro ao remover arquivo do R2: %w",
			err,
		)
	}

	return nil
}

func (s *R2Storage) PresignGetURL(
	ctx context.Context,
	key string,
	expires time.Duration,
) (string, error) {
	if key == "" {
		return "", fmt.Errorf("storage key não pode ser vazia")
	}

	if expires <= 0 || expires > 7*24*time.Hour {
		return "", fmt.Errorf("prazo de expiração inválido")
	}

	presignClient := s3.NewPresignClient(s.client)

	result, err := presignClient.PresignGetObject(
		ctx,
		&s3.GetObjectInput{
			Bucket: aws.String(s.bucket),
			Key:    aws.String(key),
		},
		func(options *s3.PresignOptions) {
			options.Expires = expires
		},
	)

	if err != nil {
		return "", fmt.Errorf(
			"erro ao gerar URL temporária do R2: %w",
			err,
		)
	}

	return result.URL, nil
}
