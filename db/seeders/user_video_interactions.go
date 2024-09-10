package seeders

import (
	"encoding/csv"
	"fmt"
	"live/auth/models"
	videomodels "live/videohub/models"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"gorm.io/gorm"
)

var eventTypes = []string{"play", "pause", "complete", "like", "dislike"}

func SeedUserVideoInteractions(db *gorm.DB) {
	var users []models.User
	var videos []videomodels.Video

	// 全てのユーザーとビデオを取得
	if err := db.Find(&users).Error; err != nil {
		log.Printf("Failed to retrieve users: %v", err)
		return
	}

	if err := db.Find(&videos).Error; err != nil {
		log.Printf("Failed to retrieve videos: %v", err)
		return
	}

	// CSVファイル用のデータを保持するスライス
	var records [][]string

	for i := 1; i <= 1000; i++ {
		// ランダムなユーザーとビデオを選択
		user := users[rand.Intn(len(users))]
		video := videos[rand.Intn(len(videos))]
		eventType := eventTypes[rand.Intn(len(eventTypes))]
		timestamp := time.Now().Unix()

		// `UserVideoInteraction`を生成
		interaction := videomodels.UserVideoInteraction{
			UserID:    user.ID,
			VideoID:   video.ID,
			EventType: eventType,
			CreatedAt: time.Now(),
		}

		// データベースにインタラクションを保存
		if err := db.Create(&interaction).Error; err != nil {
			log.Printf("Failed to create user video interaction: %v", err)
		} else {
			fmt.Printf("Created user video interaction: USER_ID=%d, ITEM_ID=%d, EVENT_TYPE=%s\n", user.ID, video.ID, eventType)

			// CSVに書き出すレコードを追加
			record := []string{
				fmt.Sprint(interaction.UserID),  // USER_ID
				fmt.Sprint(interaction.VideoID), // ITEM_ID
				interaction.EventType,           // EVENT_TYPE
				fmt.Sprint(0.0),                 // EVENT_VALUE: optional なので空の値を設定
				fmt.Sprint(timestamp),           // TIMESTAMP: UNIXタイム形式
			}
			records = append(records, record)
		}
	}

	// CSVファイルにエクスポート
	headers := []string{"USER_ID", "ITEM_ID", "EVENT_TYPE", "EVENT_VALUE", "TIMESTAMP"}
	err := createUserVideoInteractionsCSV("db/learningdata/user_video_interactions.csv", headers, records)
	if err != nil {
		fmt.Printf("Failed to export user video interactions to CSV: %v\n", err)
	}
}

// CSVファイル作成関数
func createUserVideoInteractionsCSV(filePath string, headers []string, records [][]string) error {
	dir := filepath.Dir(filePath)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			return fmt.Errorf("failed to create directory: %v", err)
		}
	}

	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// ヘッダーを書き込む
	err = writer.Write(headers)
	if err != nil {
		return fmt.Errorf("failed to write CSV header: %v", err)
	}

	// レコードを書き込む
	for _, record := range records {
		err := writer.Write(record)
		if err != nil {
			return fmt.Errorf("failed to write record: %v", err)
		}
	}

	return nil
}
