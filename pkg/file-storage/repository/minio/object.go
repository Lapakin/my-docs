package minio

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type ObjectStorageRepository struct {
	s3 *s3.Client
}

func NewObjectStorageRepository(s3 *s3.Client) *ObjectStorageRepository {
	return &ObjectStorageRepository{
		s3: s3,
	}
}

func (r *ObjectStorageRepository) Upload(ctx context.Context, path string, content io.Reader, size int64, contentType string) error {
	_, err := r.s3.PutObject(ctx, &s3.PutObjectInput{
		Bucket:        aws.String("documents"),
		Key:           aws.String(path),
		Body:          content,
		ContentLength: aws.Int64(size),
		ContentType:   aws.String(contentType),
	})

	if err != nil {
		return fmt.Errorf("failed to upload object: %w", err)
	}

	return nil
}

func (r *ObjectStorageRepository) Download(ctx context.Context, path string) (io.ReadCloser, error) {
	result, err := r.s3.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String("documents"),
		Key:    aws.String(path),
	})

	if err != nil {
		return nil, fmt.Errorf("failed to download object: %w", err)
	}

	return result.Body, nil
}

func (r *ObjectStorageRepository) Delete(ctx context.Context, path string) error {
	_, err := r.s3.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String("documents"),
		Key:    aws.String(path),
	})

	if err != nil {
		return fmt.Errorf("failed to delete object: %w", err)
	}

	return nil
}

func (r *ObjectStorageRepository) GetPresignedURL(ctx context.Context, path string) (string, error) {
	presignClient := s3.NewPresignClient(r.s3)

	request, err := presignClient.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String("documents"),
		Key:    aws.String(path),
	}, s3.WithPresignExpires(15*time.Minute))

	if err != nil {
		return "", fmt.Errorf("failed to create presigned URL: %w", err)
	}

	return request.URL, nil
}

func (r *ObjectStorageRepository) Copy(ctx context.Context, sourcePath, destPath string) error {
	copySource := fmt.Sprintf("%s/%s", "documents", sourcePath)

	_, err := r.s3.CopyObject(ctx, &s3.CopyObjectInput{
		Bucket:     aws.String("documents"),
		CopySource: aws.String(copySource),
		Key:        aws.String(destPath),
	})

	if err != nil {
		return fmt.Errorf("failed to copy object: %w", err)
	}

	return nil
}
