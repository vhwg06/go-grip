package webapi

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

type LocalStorage struct {
	baseURL string
}

func NewLocalStorage(baseURL string) *LocalStorage {
	return &LocalStorage{
		baseURL: strings.TrimSuffix(baseURL, "/"),
	}
}

func (s *LocalStorage) GeneratePresignedURL(ctx context.Context, fileName string, contentType string) (uploadURL string, publicURL string, fileID string, err error) {
	fileID = uuid.New().String()
	ext := guessExtension(contentType, fileName)
	newFileName := fileID + ext

	uploadURL = fmt.Sprintf("%s/v1/media/simulate-upload/%s", s.baseURL, newFileName)
	publicURL = fmt.Sprintf("%s/static/uploads/%s", s.baseURL, newFileName)
	return uploadURL, publicURL, fileID, nil
}

func (s *LocalStorage) Delete(ctx context.Context, key string) error {
	// Local storage is simulated, nothing to physically delete.
	return nil
}
