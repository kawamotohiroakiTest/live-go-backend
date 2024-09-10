package seeders

import (
	"encoding/csv"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"live/auth/models"

	"gorm.io/gorm"
)

// 日本語の名前のサンプル
var names = []string{
	"田中太郎", "鈴木次郎", "佐藤三郎", "高橋四郎", "伊藤五郎",
	"山本六郎", "中村七郎", "小林八郎", "加藤九郎", "渡辺十郎",
}

// ランダムな日付を今年の1月から今日までの間で生成する関数
func randomDate() time.Time {
	start := time.Date(time.Now().Year(), time.January, 1, 0, 0, 0, 0, time.Local)
	end := time.Now()
	delta := end.Unix() - start.Unix()

	sec := rand.Int63n(delta) + start.Unix()
	return time.Unix(sec, 0)
}

// ユーザーシード関数
func SeedUsers(db *gorm.DB) {
	// パスワードのハッシュを環境変数から取得
	passwordHash := os.Getenv("PASSWORD_HASH")

	if passwordHash == "" {
		log.Fatal("PASSWORD_HASH is not set in environment variables")
	}

	var records [][]string

	// 100件のユーザーを作成
	for i := 1; i <= 100; i++ {
		// メールアドレスと名前を作成
		mail := fmt.Sprintf("test%d@gmail.com", i)
		name := names[rand.Intn(len(names))]

		// ランダムなログイン日時を生成
		lastLoginAt := randomDate()

		// ユーザー作成 (models.User を使用)
		user := models.User{
			Name:        name,
			Mail:        mail,
			Pass:        passwordHash,
			LastLoginAt: lastLoginAt,
			CreatedAt:   time.Now(),
			ModifiedAt:  time.Now(),
		}

		// バリデーションチェック
		if err := user.Validate(); err != nil {
			fmt.Printf("Validation failed for user %s: %v\n", mail, err)
			continue
		}

		// 既存ユーザーを確認して、存在しない場合のみ作成
		if err := db.Where("mail = ?", mail).FirstOrCreate(&user).Error; err != nil {
			fmt.Printf("Failed to create or find user: %v\n", err)
		} else {
			fmt.Printf("Created or found user: %s\n", mail)
		}

		// CSVファイル用にレコードを追加
		record := []string{
			fmt.Sprint(user.ID),
			fmt.Sprint(user.LastLoginAt.Unix()), // タイムスタンプ形式で出力
		}
		records = append(records, record)
	}

	// CSVにエクスポート
	headers := []string{"USER_ID", "LAST_LOGIN"}
	err := createUsersCSV("db/learningdata/users.csv", headers, records)
	if err != nil {
		fmt.Printf("Failed to export users to CSV: %v\n", err)
	}
}

// CSVファイル作成関数
func createUsersCSV(filePath string, headers []string, records [][]string) error {
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

	err = writer.Write(headers)
	if err != nil {
		return fmt.Errorf("failed to write CSV header: %v", err)
	}

	for _, record := range records {
		err := writer.Write(record)
		if err != nil {
			return fmt.Errorf("failed to write record: %v", err)
		}
	}

	return nil
}
