package services

import (
	"context"
	"fmt"
	"live/common"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type StorageService struct {
	Client *minio.Client
	Bucket string
}

func NewStorageService() (*StorageService, error) {
	endpoint := os.Getenv("STORAGE_ENDPOINT")
	accessKeyID := os.Getenv("STORAGE_ACCESS_KEY")
	secretAccessKey := os.Getenv("STORAGE_SECRET_KEY")
	useSSL := false

	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		common.LogVideoUploadError(fmt.Errorf("MinIOクライアントの初期化に失敗しました: %w", err))
		return nil, err
	}

	bucketName := os.Getenv("STORAGE_BUCKET")

	exists, err := client.BucketExists(context.Background(), bucketName)
	if err != nil {
		common.LogVideoUploadError(fmt.Errorf("バケットの存在確認に失敗しました: %w", err))
		return nil, err
	}
	if !exists {
		err = client.MakeBucket(context.Background(), bucketName, minio.MakeBucketOptions{})
		if err != nil {
			common.LogVideoUploadError(fmt.Errorf("バケットの作成に失敗しました: %w", err))
			return nil, err
		}
	}

	return &StorageService{
		Client: client,
		Bucket: bucketName,
	}, nil
}

func (s *StorageService) UploadFile(file multipart.File, fileHeader *multipart.FileHeader) (string, error) {
	objectName := fileHeader.Filename

	ext := filepath.Ext(objectName)
	validExtensions := map[string]bool{
		".mp4":  true,
		".avi":  true,
		".mov":  true,
		".wmv":  true,
		".flv":  true,
		".mkv":  true,
		".webm": true,
		".mpeg": true,
		".mpg":  true,
		".3gp":  true,
		".m4v":  true,
	}

	if !validExtensions[ext] {
		err := fmt.Errorf("無効なファイル拡張子: %s", ext)
		common.LogVideoUploadError(err) // ここでログを記録
		return "", err
	}

	contentType := fileHeader.Header.Get("Content-Type")

	// ファイルをMinIOにアップロード
	_, err := s.Client.PutObject(context.Background(), s.Bucket, objectName, file, fileHeader.Size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		common.LogVideoUploadError(fmt.Errorf("ファイルのアップロードに失敗しました: %w", err))
		return "", err
	}

	fileURL := fmt.Sprintf("%s/%s/%s", os.Getenv("STORAGE_ENDPOINT"), s.Bucket, objectName)
	return fileURL, nil
}