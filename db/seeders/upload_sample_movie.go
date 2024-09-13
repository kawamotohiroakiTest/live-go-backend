package seeders

import (
	"fmt"
	"io/ioutil"
	"live/videoupload/services"
	"mime/multipart"
	"os"
	"path/filepath"
)

func UploadAllMovies() error {
	fmt.Println("動画ファイルのアップロードを開始します...")
	// ENV_MODEの取得
	envMode := os.Getenv("ENV_MODE")

	var storageService *services.StorageService
	var err error

	// ストレージサービスの初期化
	if envMode == "local" {
		fmt.Println("minio...")
		storageService, err = services.InitMinioService()
	} else {
		fmt.Println("s3...")
		storageService, err = services.NewStorageService()
	}

	if err != nil {
		return fmt.Errorf("ストレージサービスの初期化に失敗しました: %v", err)
	}

	moviesDir := filepath.Join("db", "movies")

	// ディレクトリ内の全てのファイルを取得
	files, err := ioutil.ReadDir(moviesDir)
	if err != nil {
		fmt.Println("err", err)
		return fmt.Errorf("ディレクトリの読み込みに失敗しました: %v", err)
	}
	fmt.Println("files", files)

	// 各ファイルをアップロード
	for _, fileInfo := range files {
		// ファイルが動画ファイル (.mp4) であるかを確認
		if filepath.Ext(fileInfo.Name()) != ".mp4" {
			continue // .mp4 以外のファイルは無視
		}

		// ファイルパスを作成
		moviePath := filepath.Join(moviesDir, fileInfo.Name())

		// ファイルを開く
		file, err := os.Open(moviePath)
		if err != nil {
			fmt.Printf("ファイル %s の取得に失敗しました: %v\n", fileInfo.Name(), err)
			continue
		}
		defer file.Close()

		// ファイルヘッダーを作成
		fileHeader := &multipart.FileHeader{
			Filename: fileInfo.Name(),
			Size:     fileInfo.Size(),
		}

		// 動画ファイルをアップロード
		fileURL, err := storageService.UploadFile(file, fileHeader, true)
		if err != nil {
			fmt.Printf("ファイル %s のアップロードに失敗しました: %v\n", fileInfo.Name(), err)
			continue
		}

		// アップロード結果をログ出力
		fmt.Printf("ファイル %s のアップロードに成功しました: %s\n", fileInfo.Name(), fileURL)
	}

	return nil
}
