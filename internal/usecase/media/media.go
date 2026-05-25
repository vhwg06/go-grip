package media

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/evrone/go-clean-template/internal/entity"
	"github.com/evrone/go-clean-template/internal/repo"
	"github.com/google/uuid"
)

type UseCase struct {
	repo     repo.MediaRepo
	maxBytes int64
}

func New(r repo.MediaRepo, maxBytes int64) *UseCase {
	if maxBytes <= 0 {
		maxBytes = entity.MaxMediaUploadBytes
	}
	return &UseCase{repo: r, maxBytes: maxBytes}
}

func (uc *UseCase) Store(ctx context.Context, media entity.MediaAsset) (entity.MediaAsset, error) {
	if media.SizeBytes > uc.maxBytes || !entity.IsAllowedMediaMimeType(media.MimeType) {
		return entity.MediaAsset{}, entity.ErrInvalidInput
	}
	media.ID = uuid.New().String()
	media.CreatedAt = time.Now().UTC()
	if err := uc.repo.Store(ctx, &media); err != nil {
		return entity.MediaAsset{}, err
	}
	return media, nil
}

func (uc *UseCase) List(ctx context.Context, page entity.Pagination) ([]entity.MediaAsset, int, error) {
	return uc.repo.List(ctx, page.Normalize())
}

func (uc *UseCase) Delete(ctx context.Context, id string) error {
	return uc.repo.Delete(ctx, id)
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

func (uc *UseCase) GeneratePresignedURL(ctx context.Context, fileName string, contentType string) (uploadURL string, publicURL string, fileID string, err error) {
	if !entity.IsAllowedMediaMimeType(contentType) {
		return "", "", "", entity.ErrInvalidInput
	}

	accountID := os.Getenv("R2_ACCOUNT_ID")
	accessKey := os.Getenv("R2_ACCESS_KEY_ID")
	secretKey := os.Getenv("R2_SECRET_ACCESS_KEY")
	bucketName := os.Getenv("R2_BUCKET_NAME")
	publicURL = os.Getenv("R2_PUBLIC_URL")

	fileID = uuid.New().String()
	ext := filepath.Ext(fileName)
	if ext == "" {
		switch contentType {
		case "image/jpeg":
			ext = ".jpg"
		case "image/png":
			ext = ".png"
		case "image/webp":
			ext = ".webp"
		default:
			ext = ".bin"
		}
	}
	newFileName := fileID + ext

	if accountID == "" || accessKey == "" || secretKey == "" || bucketName == "" {
		// Môi trường Local Dev giả lập
		uploadURL = fmt.Sprintf("http://localhost:8080/api/media/simulate-upload/%s", newFileName)
		if publicURL == "" {
			publicURL = fmt.Sprintf("http://localhost:8080/static/uploads/%s", newFileName)
		} else {
			publicURL = fmt.Sprintf("%s/%s", strings.TrimSuffix(publicURL, "/"), newFileName)
		}
		return uploadURL, publicURL, fileID, nil
	}

	resolver := &r2EndpointResolver{
		accountID: accountID,
	}

	client := s3.New(s3.Options{
		Credentials:      credentials.NewStaticCredentialsProvider(accessKey, secretKey, ""),
		Region:           "auto",
		EndpointResolver: resolver,
	})

	presignClient := s3.NewPresignClient(client)
	presignedReq, err := presignClient.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(bucketName),
		Key:         aws.String(newFileName),
		ContentType: aws.String(contentType),
	}, s3.WithPresignExpires(15*time.Minute))
	if err != nil {
		return "", "", "", fmt.Errorf("failed to generate presigned URL: %w", err)
	}

	uploadURL = presignedReq.URL
	if publicURL == "" {
		publicURL = fmt.Sprintf("https://%s.r2.cloudflarestorage.com/%s/%s", accountID, bucketName, newFileName)
	} else {
		publicURL = fmt.Sprintf("%s/%s", strings.TrimSuffix(publicURL, "/"), newFileName)
	}

	return uploadURL, publicURL, fileID, nil
}
