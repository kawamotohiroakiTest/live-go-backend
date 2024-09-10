package models

import (
	"live/common"
	"time"

	"gorm.io/gorm"
)

type Video struct {
	ID          uint        `gorm:"primary_key"`
	UserID      uint        `gorm:"not null"`                                              // 外部キー usersテーブルのID
	Title       string      `gorm:"type:varchar(255);not null"`                            // 動画のタイトル
	Description string      `gorm:"type:text"`                                             // 説明
	ViewCount   uint        `gorm:"default:0;not null"`                                    // 視聴回数
	Rating      float64     `gorm:"type:decimal(3,2);default:0.00"`                        // 評価
	Genre       string      `gorm:"type:varchar(255);not null"`                            // ジャンル
	PostedAt    time.Time   `gorm:"default:CURRENT_TIMESTAMP;not null"`                    // 投稿日時
	Created     time.Time   `gorm:"default:CURRENT_TIMESTAMP"`                             // 作成日時
	Modified    time.Time   `gorm:"default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"` // 更新日時
	Deleted     *time.Time  `gorm:"default:NULL"`                                          // 削除日時
	Files       []VideoFile `gorm:"foreignKey:VideoID"`                                    // 動画ファイルとのリレーション
}

type VideoFile struct {
	ID            uint       `gorm:"primary_key"`                                                                       // ファイルのID
	VideoID       uint       `gorm:"not null"`                                                                          // 外部キー、videosテーブルのID
	FilePath      string     `gorm:"type:varchar(255);not null"`                                                        // 動画ファイルパス
	ThumbnailPath string     `gorm:"type:varchar(255)"`                                                                 // サムネイルパス
	Duration      uint       `gorm:"not null"`                                                                          // 動画の再生時間 (秒単位)
	FileSize      uint64     `gorm:"not null"`                                                                          // ファイルサイズ (バイト単位)
	Format        string     `gorm:"type:varchar(50);not null"`                                                         // ファイル形式 (例: mp4)
	Status        string     `gorm:"type:enum('pending','processing','completed','failed');default:'pending';not null"` // ステータス
	Created       time.Time  `gorm:"default:CURRENT_TIMESTAMP"`                                                         // 作成日時
	Modified      time.Time  `gorm:"default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`                             // 更新日時
	Deleted       *time.Time `gorm:"default:NULL"`                                                                      // 削除日時
}

func SaveVideo(userID uint, title, description string) (*Video, error) {
	video := Video{
		UserID:      userID,
		Title:       title,
		Description: description,
	}

	if err := common.DB.Create(&video).Error; err != nil {
		return nil, err
	}

	return &video, nil
}

func SaveVideoFile(videoID uint, filePath, thumbnailPath string, duration uint, fileSize uint64, format string) (*VideoFile, error) {
	videoFile := VideoFile{
		VideoID:       videoID,
		FilePath:      filePath,
		ThumbnailPath: thumbnailPath,
		Duration:      duration,
		FileSize:      fileSize,
		Format:        format,
	}

	if err := common.DB.Create(&videoFile).Error; err != nil {
		return nil, err
	}

	return &videoFile, nil
}

// トランザクションを使用して動画情報を保存する関数
func SaveVideoWithTransaction(tx *gorm.DB, userID uint, title, description string) (*Video, error) {
	video := Video{
		UserID:      userID,
		Title:       title,
		Description: description,
	}

	if err := tx.Create(&video).Error; err != nil {
		return nil, err
	}

	return &video, nil
}

// トランザクションを使用して動画ファイル情報を保存する関数
func SaveVideoFileWithTransaction(tx *gorm.DB, videoID uint, filePath string, duration uint, fileSize uint64, format string) (*VideoFile, error) {
	videoFile := VideoFile{
		VideoID:  videoID,
		FilePath: filePath,
		Duration: duration,
		FileSize: fileSize,
		Format:   format,
	}

	if err := tx.Create(&videoFile).Error; err != nil {
		return nil, err
	}

	return &videoFile, nil
}
