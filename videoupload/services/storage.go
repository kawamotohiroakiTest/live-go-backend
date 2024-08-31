package services

import (
	"bytes"
	"fmt"
	"live/common"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
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
		region = "ap-northeast-1" // デフォルトのリージョン
	}

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})
	if err != nil {
		common.LogVideoUploadError(fmt.Errorf("AWSセッションの初期化に失敗しました: %w", err))
		return nil, err
	}

	client := s3.New(sess)

	bucketName := os.Getenv("STORAGE_BUCKET")
	common.LogVideoUploadError(fmt.Errorf("Using Bucket=%s", bucketName))

	return &StorageService{
		Client: client,
		Bucket: bucketName,
	}, nil
}

// 動画ファイルをアップロードするメソッド
func (s *StorageService) UploadFile(file multipart.File, fileHeader *multipart.FileHeader) (string, error) {
	// ファイル名を `movies/` プレフィックスにする
	objectName := "movies/" + common.GenerateUniqueFileName(fileHeader.Filename)

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
		common.LogVideoUploadError(err)
		return "", err
	}

	contentType := fileHeader.Header.Get("Content-Type")

	// ファイルの内容をバイトスライスに変換
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(file)
	if err != nil {
		common.LogVideoUploadError(fmt.Errorf("ファイルの読み込みに失敗しました: %w", err))
		return "", err
	}
	fileBytes := buf.Bytes()

	// ファイルをS3にアップロード
	_, err = s.Client.PutObject(&s3.PutObjectInput{
		Bucket:      aws.String(s.Bucket),
		Key:         aws.String(objectName),
		Body:        bytes.NewReader(fileBytes),
		ContentType: aws.String(contentType),
	})
	if err != nil {
		// エラーの種類と詳細をログに出力
		common.LogVideoUploadError(fmt.Errorf(
			"ファイルのアップロードに失敗しました: %w | Bucket: %s, Key: %s, Content-Type: %s",
			err,
			s.Bucket,
			objectName,
			contentType,
		))
		return "", err
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
