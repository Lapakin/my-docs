package minio

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/lapotkin/file-storage/internal/adapter/object"

	transport "github.com/aws/smithy-go/endpoints"
)

const region = "us-east-1"

type minioEndpointResolver struct {
	URL *url.URL
}

func (r *minioEndpointResolver) ResolveEndpoint(_ context.Context, params s3.EndpointParameters) (transport.Endpoint, error) {
	u := *r.URL
	u.Path += "/" + *params.Bucket
	return transport.Endpoint{URI: u}, nil
}

func NewMinIOClient(cfg *object.Config) (*s3.Client, error) {
	endpoint := cfg.Endpoint
	if !strings.HasPrefix(endpoint, "http://") && !strings.HasPrefix(endpoint, "https://") {
		protocol := "http"
		if cfg.UseSSL {
			protocol = "https"
		}
		endpoint = fmt.Sprintf("%s://%s", protocol, endpoint)
	}

	endpointURL, err := url.Parse(endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to parse endpoint URL: %w", err)
	}
	client := s3.New(s3.Options{
		Credentials: aws.CredentialsProviderFunc(func(_ context.Context) (aws.Credentials, error) {
			return aws.Credentials{
				AccessKeyID:     cfg.AccessKey,
				SecretAccessKey: cfg.SecretKey,
			}, nil
		}),
		Region:             region,
		EndpointResolverV2: &minioEndpointResolver{URL: endpointURL},
		UsePathStyle:       true,
	})

	if cfg.BucketName != "" {
		_, err = client.HeadBucket(context.Background(), &s3.HeadBucketInput{
			Bucket: aws.String(cfg.BucketName),
		})
		if err != nil {
			_, err = client.CreateBucket(context.Background(), &s3.CreateBucketInput{
				Bucket: aws.String(cfg.BucketName),
			})
			if err != nil {
				return nil, fmt.Errorf("failed to create bucket %s: %w", cfg.BucketName, err)
			}
		}
	}

	return client, nil
}
