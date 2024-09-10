package seeders

import (
	"encoding/csv"
	"fmt"
	"live/auth/models"                 // usersテーブルのエイリアスをauthに設定
	videomodels "live/videohub/models" // videoモデルのエイリアスを設定
	"live/videoupload/services"
	"math/rand"
	"mime/multipart"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// 動画ジャンルのサンプル
var genres = []string{
	"アクション", "コメディ", "ドキュメンタリー", "ドラマ", "ホラー",
	"SF", "アニメ", "音楽", "スポーツ", "ニュース",
}

// サンプル動画を生成する関数
func generateSampleVideo(outputPath string) error {
	// 出力ディレクトリが存在しない場合は作成
	outputDir := filepath.Dir(outputPath)
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}
	}

	// ffmpeg コマンドを使用して動画を生成（テキスト描画は省略）
	cmd := exec.Command("ffmpeg", "-f", "lavfi", "-i", fmt.Sprintf("color=c=%s:s=1280x720:d=5", randomColor()),
		"-loglevel", "quiet",
		"-c:v", "libx264", "-t", "5", "-pix_fmt", "yuv420p", outputPath)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// コマンドを実行
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to generate video: %w", err)
	}
	return nil
}

// ファイルサイズを取得する関数
func getFileSize(filePath string) int64 {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return 0
	}
	return fileInfo.Size()
}

// ファイルをmultipart.Fileとして読み込む関数
func createMultipartFile(filePath string) (multipart.File, *multipart.FileHeader, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, nil, err
	}

	fileHeader := &multipart.FileHeader{
		Filename: filepath.Base(filePath),
		Size:     getFileSize(filePath),
		Header:   make(map[string][]string),
	}
	fileHeader.Header.Set("Content-Type", "video/mp4")

	return file, fileHeader, nil
}

// ランダムなユーザーをデータベースから取得する関数
func getRandomUser(db *gorm.DB) (*models.User, error) {
	var user models.User
	if err := db.Order("RAND()").Limit(1).Find(&user).Error; err != nil {
		return nil, fmt.Errorf("failed to get random user: %w", err)
	}
	return &user, nil
}

// シーダー関数の例
func SeedVideos(db *gorm.DB) {
	var records [][]string // CSV用のレコード

	for i := 1; i <= 100; i++ {
		// データベースからランダムにユーザーを取得
		user, err := getRandomUser(db)
		if err != nil {
			fmt.Printf("Failed to get random user: %v\n", err)
			continue
		}

		genre := genres[rand.Intn(len(genres))]
		title := fmt.Sprintf("%s映画%d", genre, i)
		description := fmt.Sprintf("%s映画の説明%d", genre, i)

		// サンプル動画ファイルを生成
		videoFilePath := fmt.Sprintf("movies/%s.mp4", uuid.New().String())
		err = generateSampleVideo(videoFilePath)
		if err != nil {
			fmt.Printf("Failed to generate video: %v\n", err)
			continue
		}

		// 動画ファイルを multipart.File として読み込む
		file, fileHeader, err := createMultipartFile(videoFilePath)
		if err != nil {
			fmt.Printf("Failed to create multipart file: %v\n", err)
			continue
		}
		defer file.Close()

		// 動画ファイルのアップロード
		fileURL, err := services.UploadVideoFile(user.ID, title, description, file, fileHeader, "300")
		if err != nil {
			fmt.Printf("Failed to upload video: %v\n", err)
			continue
		}

		// 動画ファイル情報をデータベースに保存
		video := videomodels.Video{
			UserID:      user.ID, // ランダムに取得されたユーザーのIDを使用
			Title:       title,
			Description: description,
			ViewCount:   uint(i),
			Rating:      0.00,
			Genre:       genre,
			PostedAt:    time.Now(),
			Created:     time.Now(),
			Modified:    time.Now(),
		}

		if err := db.Create(&video).Error; err != nil {
			fmt.Printf("Failed to create video: %v\n", err)
			continue
		}

		// videoFileデータを保存
		videoFile := videomodels.VideoFile{
			VideoID:  video.ID,
			FilePath: fileURL,
			Format:   "mp4",
			Status:   "pending",
			Created:  time.Now(),
			Modified: time.Now(),
		}

		if err := db.Create(&videoFile).Error; err != nil {
			fmt.Printf("Failed to create video file: %v\n", err)
			continue
		}

		// CSV用のレコードを作成
		record := []string{
			fmt.Sprint(video.ID),              // ITEM_ID
			video.Genre,                       // GENRES
			fmt.Sprint(video.Created.Unix()),  // CREATION_TIMESTAMP (UNIXタイムスタンプ)
			fmt.Sprint(video.ViewCount),       // VIEW_COUNT
			fmt.Sprintf("%.2f", video.Rating), // RATING
		}
		records = append(records, record)
	}

	// moviesディレクトリを削除
	err := clearMoviesDirectory("movies")
	if err != nil {
		fmt.Printf("Failed to clear movies directory: %v\n", err)
	}

	// CSVファイルの作成
	headers := []string{"ITEM_ID", "GENRES", "CREATION_TIMESTAMP", "VIEW_COUNT", "RATING"}
	err = createVideosCSV("db/learningdata/videos.csv", headers, records)
	if err != nil {
		fmt.Printf("Failed to create videos CSV: %v\n", err)
	}
}

// movies ディレクトリを削除する関数
func clearMoviesDirectory(dirPath string) error {
	err := os.RemoveAll(dirPath)
	if err != nil {
		return fmt.Errorf("Failed to remove movies directory: %v", err)
	}
	return nil
}

// CSVファイルを作成する関数
func createVideosCSV(filePath string, headers []string, records [][]string) error {
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

func randomColor() string {
	colors := []string{"red", "blue", "green", "yellow", "purple", "orange", "pink", "cyan"}
	return colors[rand.Intn(len(colors))]
}
