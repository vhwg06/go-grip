package webapi

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
)

type R2Storage struct {
	client        *s3.Client
	presignClient *s3.PresignClient
	bucketName    string
	publicURL     string
	accountID     string
}

type r2EndpointResolver struct {
	accountID string
}

func (r *r2EndpointResolver) ResolveEndpoint(region string, options s3.EndpointResolverOptions) (aws.Endpoint, error) {
	return aws.Endpoint{
		URL:           fmt.Sprintf("https://%s.r2.cloudflarestorage.com", r.accountID),
		SigningRegion: "auto",
	}, nil
}

// NewR2Storage creates a singleton S3/R2 client initialized once on startup.
func NewR2Storage(accountID, accessKey, secretKey, bucketName, publicURL string) *R2Storage {
	resolver := &r2EndpointResolver{
		accountID: accountID,
	}

	client := s3.New(s3.Options{
		Credentials:      credentials.NewStaticCredentialsProvider(accessKey, secretKey, ""),
		Region:           "auto",
		EndpointResolver: resolver,
	})

	presignClient := s3.NewPresignClient(client)

	return &R2Storage{
		client:        client,
		presignClient: presignClient,
		bucketName:    bucketName,
		publicURL:     publicURL,
		accountID:     accountID,
	}
}

func (s *R2Storage) GeneratePresignedURL(ctx context.Context, fileName string, contentType string) (uploadURL string, publicURL string, fileID string, err error) {
	fileID = uuid.New().String()
	ext := guessExtension(contentType, fileName)
	newFileName := fileID + ext

	presignedReq, err := s.presignClient.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s.bucketName),
		Key:         aws.String(newFileName),
		ContentType: aws.String(contentType),
	}, s3.WithPresignExpires(15*time.Minute))
	if err != nil {
		return "", "", "", err
	}

	uploadURL = presignedReq.URL
	if s.publicURL == "" {
		publicURL = fmt.Sprintf("https://%s.r2.cloudflarestorage.com/%s/%s", s.accountID, s.bucketName, newFileName)
	} else {
		publicURL = fmt.Sprintf("%s/%s", strings.TrimSuffix(s.publicURL, "/"), newFileName)
	}

	return uploadURL, publicURL, fileID, nil
}

func (s *R2Storage) Delete(ctx context.Context, key string) error {
	_, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(key),
	})
	return err
}
