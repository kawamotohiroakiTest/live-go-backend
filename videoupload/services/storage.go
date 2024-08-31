package services

import (
	"bytes"
	"context"
	"fmt"
	"live/common"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type StorageService struct {
	Client *s3.S3
	Bucket string
}

func NewStorageService() (*StorageService, error) {
	region := os.Getenv("AWS_REGION")
	if region == "" {
		return nil, fmt.Errorf("AWS_REGION 環境変数が設定されていません")
	}

	bucketName := os.Getenv("STORAGE_BUCKET")
	if bucketName == "" {
		return nil, fmt.Errorf("STORAGE_BUCKET 環境変数が設定されていません")
	}

	accessKey := os.Getenv("AWS_ACCESS_KEY_ID")
	if accessKey == "" {
		return nil, fmt.Errorf("AWS_ACCESS_KEY_ID 環境変数が設定されていません")
	}

	secretKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	if secretKey == "" {
		return nil, fmt.Errorf("AWS_SECRET_ACCESS_KEY 環境変数が設定されていません")
	}

	// セッションの作成
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(region),
		Credentials: credentials.NewStaticCredentials(accessKey, secretKey, ""),
		LogLevel:    aws.LogLevel(aws.LogDebugWithHTTPBody),
	})
	if err != nil {
		return nil, fmt.Errorf("AWSセッションの初期化に失敗しました: %w", err)
	}

	client := s3.New(sess)
	common.LogVideoUploadError(fmt.Errorf("Using Bucket=%s", bucketName))

	return &StorageService{
		Client: client,
		Bucket: bucketName,
	}, nil
}

func (s *StorageService) UploadFile(file multipart.File, fileHeader *multipart.FileHeader) (string, error) {
	defer file.Close() // ファイルを閉じるためのdefer

	objectName := "movies/" + common.GenerateUniqueFileName(fileHeader.Filename)

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

	ext := filepath.Ext(objectName)
	if !validExtensions[ext] {
		return "", fmt.Errorf("無効なファイル拡張子: %s", ext)
	}

	contentType := fileHeader.Header.Get("Content-Type")

	// ファイルの内容をバイトスライスに変換
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(file)
	if err != nil {
		return "", fmt.Errorf("ファイルの読み込みに失敗しました: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second) // タイムアウト設定
	defer cancel()

	// ファイルをS3にアップロード
	_, err = s.Client.PutObjectWithContext(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s.Bucket),
		Key:         aws.String(objectName),
		Body:        bytes.NewReader(buf.Bytes()),
		ContentType: aws.String(contentType),
	})
	if err != nil {
		return "", fmt.Errorf("ファイルのアップロードに失敗しました: %w | Bucket: %s, Key: %s, Content-Type: %s",
			err, s.Bucket, objectName, contentType)
	}

	fileURL := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", s.Bucket, os.Getenv("AWS_REGION"), objectName)
	return fileURL, nil
}

// サムネイルをアップロードするメソッド
func (s *StorageService) UploadThumbnailFile(file multipart.File, fileHeader *multipart.FileHeader) (string, error) {
	// ファイル名を `thumbnails/` プレフィックスにする
	objectName := "thumbnails/" + common.GenerateUniqueFileName(fileHeader.Filename)

	ext := filepath.Ext(objectName)
	validExtensions := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
		".bmp":  true,
		".tiff": true,
		".webp": true,
		".ico":  true,
		".svg":  true,
	}

	if !validExtensions[ext] {
		err := fmt.Errorf("無効なファイル拡張子: %s", ext)
		common.LogVideoUploadError(err)
		return "", err
	}

	contentType := fileHeader.Header.Get("Content-Type")

	// ファイルの内容をバイトスライスに変換
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(file)
	if err != nil {
		common.LogVideoUploadError(fmt.Errorf("サムネイルファイルの読み込みに失敗しました: %w", err))
		return "", err
	}
	fileBytes := buf.Bytes()

	// サムネイルファイルをS3にアップロード
	_, err = s.Client.PutObject(&s3.PutObjectInput{
		Bucket:      aws.String(s.Bucket),
		Key:         aws.String(objectName),
		Body:        bytes.NewReader(fileBytes),
		ContentType: aws.String(contentType),
	})
	if err != nil {
		common.LogVideoUploadError(fmt.Errorf("サムネイルのアップロードに失敗しました: %w", err))
		return "", err
	}

	fileURL := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", s.Bucket, os.Getenv("AWS_REGION"), objectName)
	return fileURL, nil
}
