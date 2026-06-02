package services

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"github.com/aidilbaihaqi/prabodrive-be/internal/config"
)

type S3Service struct {
	client  *s3.Client
	presign *s3.PresignClient
	bucket  string
	expiry  time.Duration
}

func NewS3Service(cfg config.Config) (*S3Service, error) {
	awsCfg, err := awsconfig.LoadDefaultConfig(context.Background(),
		awsconfig.WithRegion(cfg.AWS.Region),
		awsconfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			cfg.AWS.AccessKeyID,
			cfg.AWS.SecretAccessKey,
			"",
		)),
	)
	if err != nil {
		return nil, fmt.Errorf("s3: load aws config: %w", err)
	}

	client := s3.NewFromConfig(awsCfg)
	return &S3Service{
		client:  client,
		presign: s3.NewPresignClient(client),
		bucket:  cfg.S3.Bucket,
		expiry:  cfg.S3.PresignExpiry,
	}, nil
}

func (s *S3Service) GeneratePutURL(ctx context.Context, s3Key, mimeType string) (string, time.Time, error) {
	expiresAt := time.Now().Add(s.expiry)

	req, err := s.presign.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s.bucket),
		Key:         aws.String(s3Key),
		ContentType: aws.String(mimeType),
	}, s3.WithPresignExpires(s.expiry))
	if err != nil {
		return "", time.Time{}, fmt.Errorf("s3: presign put: %w", err)
	}

	return req.URL, expiresAt, nil
}

func (s *S3Service) GenerateGetURL(ctx context.Context, s3Key string, expiry time.Duration) (string, time.Time, error) {
	expiresAt := time.Now().Add(expiry)

	req, err := s.presign.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(s3Key),
	}, s3.WithPresignExpires(expiry))
	if err != nil {
		return "", time.Time{}, fmt.Errorf("s3: presign get: %w", err)
	}

	return req.URL, expiresAt, nil
}

func (s *S3Service) DeleteObject(ctx context.Context, s3Key string) error {
	_, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(s3Key),
	})
	if err != nil {
		return fmt.Errorf("s3: delete object: %w", err)
	}
	return nil
}

// S3Key returns the canonical key for a document.
func S3Key(userID, folderID, docID, filename string) string {
	folder := "root"
	if folderID != "" {
		folder = folderID
	}
	return fmt.Sprintf("%s/%s/%s_%s", userID, folder, docID, filename)
}
