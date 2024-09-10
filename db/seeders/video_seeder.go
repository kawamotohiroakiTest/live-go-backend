package seeders

import (
	"fmt"
	"live/videohub/models"
	"log"
	"math/rand"
	"path/filepath"
	"time"

	"mime/multipart"
	"os"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// 動画ジャンルのサンプル
var genres = []string{
	"アクション", "コメディ", "ドキュメンタリー", "ドラマ", "ホラー",
	"SF", "アニメ", "音楽", "スポーツ", "ニュース",
}

// シーダー関数の例
func SeedVideos(db *gorm.DB) {
	for i := 1; i <= 100; i++ {
		// ユーザーIDをランダムに取得
		userID := rand.Intn(100) + 1
		genre := genres[rand.Intn(len(genres))]
		title := fmt.Sprintf("%s映画%d", genre, i)
		description := fmt.Sprintf("%s映画の説明%d", genre, i)

		// 動画のローカルパス（ここではサンプル動画を生成するか、既存の動画を指定）
		filePath := fmt.Sprintf("movies/%s.mp4", uuid.New().String())

		video := models.Video{
			UserID:      uint(userID),
			Title:       title,
			Description: description,
			ViewCount:   uint(rand.Intn(100)), // 視聴回数をランダムに生成
			Rating:      0.00,                 // 評価は0
			Genre:       genre,
			PostedAt:    time.Now(),
			Created:     time.Now(),
			Modified:    time.Now(),
		}

		if err := db.Create(&video).Error; err != nil {
			log.Printf("Failed to create video: %v", err)
			continue
		}

		// 動画ファイル情報をデータベースに保存
		videoFile := models.VideoFile{
			VideoID:  video.ID,
			FilePath: fileURL, // MinIO上の動画のURL
			Format:   "mp4",
			Status:   "pending",
			Created:  time.Now(),
			Modified: time.Now(),
		}

		if err := db.Create(&videoFile).Error; err != nil {
			log.Printf("Failed to create video file: %v", err)
		}
	}
}

// 仮に動画ファイルを生成するための関数
func createMultipartFile(filePath string) (multipart.File, *multipart.FileHeader) {
	file, _ := os.Open(filePath)
	fileHeader := &multipart.FileHeader{
		Filename: filepath.Base(filePath),
		Size:     getFileSize(filePath),
	}
	return file, fileHeader
}

// ファイルサイズを取得する関数
func getFileSize(filePath string) int64 {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return 0
	}
	return fileInfo.Size()
}
