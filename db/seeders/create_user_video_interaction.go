package seeders

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var names = []string{
	"佐藤", "鈴木", "高橋", "田中", "伊藤", "山本", "中村", "小林", "加藤", "吉田",
	"山田", "佐々木", "松本", "井上", "木村", "林", "清水", "斎藤", "山口", "池田",
	"橋本", "阿部", "石川", "森", "遠藤", "青木", "藤田", "西村", "福田", "岡田",
	"中島", "小川", "長谷川", "村上", "近藤", "坂本", "石井", "岡本", "和田", "竹内",
	"金子", "中山", "藤井", "上田", "森田", "原田", "柴田", "酒井", "工藤", "横山",
}

var genres = []string{
	"Action", "Comedy", "Documentary", "Drama", "Horror", "Sci-Fi", "Anime", "Music", "Sports", "News",
}

type User struct {
	ID          int `gorm:"primaryKey;autoIncrement"`
	Name        string
	Mail        string
	Pass        string
	LastLoginAt time.Time
	CreatedAt   time.Time
	ModifiedAt  time.Time
}

type Video struct {
	ID          int `gorm:"primaryKey;autoIncrement"`
	UserID      int
	Title       string
	Description string
	Created     time.Time
	Modified    time.Time
	ViewCount   uint
	Rating      float64
	Genre       string
	PostedAt    time.Time
}

type VideoFile struct {
	ID       int `gorm:"primaryKey;autoIncrement"`
	VideoID  int
	FilePath string
	Format   string
	Status   string
	Created  time.Time
	Modified time.Time
}

type UserVideoInteraction struct {
	ID        int `gorm:"primaryKey;autoIncrement"`
	UserID    int
	VideoID   int
	EventType string
	CreatedAt time.Time
}

// 120人分のユーザーデータを生成
func SeedUsers(db *gorm.DB) error {

	for i := 1; i <= 120; i++ {
		name := names[rand.Intn(len(names))]
		mail := fmt.Sprintf("test%d@gmail.com", i)

		lastLoginAt := time.Now().Add(-time.Duration(rand.Intn(365*24)) * time.Hour)

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte("testtest"), bcrypt.DefaultCost)
		if err != nil {
			fmt.Printf("Error generating bcrypt hash: %v", err)
			return err
		}

		// ユーザーデータを作成
		user := User{
			Name:        name,
			Mail:        mail,
			Pass:        string(hashedPassword),
			LastLoginAt: lastLoginAt,
			CreatedAt:   time.Now(),
			ModifiedAt:  time.Now(),
		}

		// データベースに保存
		if err := db.Create(&user).Error; err != nil {
			return fmt.Errorf("Error seeding users: %v", err)
		}
	}

	fmt.Println("Users seeding completed successfully.")
	return nil
}

// 動画とそのファイルを120件生成
func SeedVideos(db *gorm.DB) error {
	genreCount := len(genres)

	for i := 1; i <= 120; i++ {
		var userID uint
		db.Raw("SELECT id FROM users ORDER BY RAND() LIMIT 1").Scan(&userID)

		genre := genres[(i-1)%genreCount]

		title := fmt.Sprintf("%sタイトル%d", genre, i)
		description := fmt.Sprintf("%s説明%d", genre, i)

		viewCount := rand.Intn(1000) + 1
		rating := rand.Float64()*4 + 1

		// ビデオデータを生成
		video := Video{
			UserID:      int(userID),
			Title:       title,
			Description: description,
			ViewCount:   uint(viewCount),
			Rating:      rating,
			Genre:       genre,
			PostedAt:    time.Now(),
			Created:     time.Now(),
			Modified:    time.Now(),
		}

		// ビデオをデータベースに保存
		if err := db.Create(&video).Error; err != nil {
			return fmt.Errorf("Error seeding videos: %v", err)
		}

		// video_files テーブルに対応するレコードを生成
		filePath := fmt.Sprintf("movies/sample_%d_%s.mp4", i, strings.ToLower(genre))
		videoFile := VideoFile{
			VideoID:  video.ID,
			FilePath: filePath,
			Format:   "mp4",
			Status:   "pending",
			Created:  time.Now(),
			Modified: time.Now(),
		}

		// video_fileをデータベースに保存
		if err := db.Create(&videoFile).Error; err != nil {
			return fmt.Errorf("Error seeding video_files: %v", err)
		}
	}

	fmt.Println("Video and VideoFile seeding completed successfully.")
	return nil
}

// ユーザーのインタラクションを1200件生成
func SeedUserVideoInteractions(db *gorm.DB) error {
	eventTypes := []string{"play", "pause", "complete", "like", "dislike"}

	var userIDs []int
	var videoIDs []int

	// ユーザーのIDを取得
	db.Model(&User{}).Pluck("id", &userIDs)
	// ビデオのIDを取得
	db.Model(&Video{}).Pluck("id", &videoIDs)

	if len(userIDs) == 0 || len(videoIDs) == 0 {
		return fmt.Errorf("No users or videos found in the database")
	}

	for i := 0; i < 1200; i++ {
		interaction := UserVideoInteraction{
			UserID:    userIDs[rand.Intn(len(userIDs))],
			VideoID:   videoIDs[rand.Intn(len(videoIDs))],
			EventType: eventTypes[rand.Intn(len(eventTypes))],
		}
		if err := db.Create(&interaction).Error; err != nil {
			return err
		}
	}
	fmt.Println("1200 user_video_interactions seeded.")
	return nil
}

// ランダムな文字列を生成する関数
func generateRandomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

// メインのSeeder関数
func SeedAll(db *gorm.DB) {
	if err := SeedUsers(db); err != nil {
		fmt.Println("Error seeding users:", err)
	}

	if err := SeedVideos(db); err != nil {
		fmt.Println("Error seeding videos and files:", err)
	}

	if err := SeedUserVideoInteractions(db); err != nil {
		fmt.Println("Error seeding user video interactions:", err)
	}
}
