package seeders

import (
	"encoding/csv"
	"fmt"
	"live/common"
	"log"
	"os"
	"path/filepath"

	"gorm.io/gorm"
)

// CSVファイルに書き込む関数
func writeCSV(filename string, headers []string, rows [][]string) error {
	// 現在の作業ディレクトリを取得
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("現在のディレクトリの取得に失敗しました: %v", err)
	}

	// 相対パスを基にCSVのパスを作成
	relativePath := filepath.Join(cwd, "db/csv", filename)

	// ディレクトリが存在しない場合は作成
	dir := filepath.Dir(relativePath)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			return fmt.Errorf("ディレクトリ作成に失敗しました: %v", err)
		}
	}

	// CSVファイルの作成
	file, err := os.Create(relativePath)
	if err != nil {
		return fmt.Errorf("ファイル作成に失敗しました: %v", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// ヘッダーを書き込む
	if err := writer.Write(headers); err != nil {
		return fmt.Errorf("ヘッダー書き込みに失敗しました: %v", err)
	}

	// 各レコードを書き込む
	for _, row := range rows {
		if err := writer.Write(row); err != nil {
			return fmt.Errorf("レコード書き込みに失敗しました: %v", err)
		}
	}

	return nil
}

// ユーザーデータを取得し、CSVに書き込む
func exportUsersToCSV(db *gorm.DB) error {
	var users []User
	if err := db.Find(&users).Error; err != nil {
		return err
	}

	// CSV形式のデータに変換
	var records [][]string
	for _, user := range users {
		record := []string{
			fmt.Sprintf("user_%d", user.ID),
			fmt.Sprintf("%d", user.LastLoginAt.Unix()),
		}
		records = append(records, record)
	}

	// CSV書き込み
	headers := []string{"USER_ID", "LAST_LOGIN"}
	return writeCSV("../csv/users.csv", headers, records)
}

// 動画データを取得し、CSVに書き込む
func exportVideosToCSV(db *gorm.DB) error {
	var videos []Video
	if err := db.Find(&videos).Error; err != nil {
		return err
	}

	// CSV形式のデータに変換
	var records [][]string
	for _, video := range videos {
		record := []string{
			fmt.Sprintf("video_%d", video.ID),
			video.Title,
			video.Genre,
			fmt.Sprintf("%d", video.Created.Unix()),
		}
		records = append(records, record)
	}

	// CSV書き込み
	headers := []string{"ITEM_ID", "TITLE", "GENRES", "CREATION_TIMESTAMP"}
	return writeCSV("../csv/videos.csv", headers, records)
}

// ユーザーと動画のインタラクションデータを取得し、CSVに書き込む
func exportInteractionsToCSV(db *gorm.DB) error {
	var interactions []UserVideoInteraction
	if err := db.Find(&interactions).Error; err != nil {
		return err
	}

	// CSV形式のデータに変換
	var records [][]string
	for _, interaction := range interactions {
		record := []string{
			fmt.Sprintf("user_%d", interaction.UserID),
			fmt.Sprintf("video_%d", interaction.VideoID),
			interaction.EventType,
			fmt.Sprintf("%d", interaction.CreatedAt.Unix()),
		}
		records = append(records, record)
	}

	// CSV書き込み
	headers := []string{"USER_ID", "ITEM_ID", "EVENT_TYPE", "TIMESTAMP"}
	return writeCSV("../csv/interactions.csv", headers, records)
}

func CreateCSV() {
	// データベース初期化
	dbConn, err := common.InitDB()
	if err != nil {
		log.Fatalf("Error initializing the database: %v", err)
	}

	// 各データをCSVにエクスポート
	if err := exportUsersToCSV(dbConn); err != nil {
		fmt.Printf("ユーザーエクスポート失敗: %v\n", err)
	}
	if err := exportVideosToCSV(dbConn); err != nil {
		fmt.Printf("動画エクスポート失敗: %v\n", err)
	}
	if err := exportInteractionsToCSV(dbConn); err != nil {
		fmt.Printf("インタラクションエクスポート失敗: %v\n", err)
	}
}
