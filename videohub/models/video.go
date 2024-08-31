package models

import (
	"live/common"
	"time"
)

type Video struct {
	ID          uint        `gorm:"primary_key"`
	Title       string      `gorm:"type:varchar(255);not null"`
	Description string      `gorm:"type:text"`
	Created     time.Time   `gorm:"default:CURRENT_TIMESTAMP"`
	Modified    time.Time   `gorm:"default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
	Deleted     *time.Time  `gorm:"default:NULL"`
	Files       []VideoFile `gorm:"foreignKey:VideoID"` // ここで動画ファイルとのリレーションを設定
}

type VideoFile struct {
	ID            uint       `gorm:"primary_key"`
	VideoID       uint       `gorm:"not null"`
	FilePath      string     `gorm:"type:varchar(255);not null"`
	ThumbnailPath string     `gorm:"type:varchar(255)"`
	Created       time.Time  `gorm:"default:CURRENT_TIMESTAMP"`
	Modified      time.Time  `gorm:"default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
	Deleted       *time.Time `gorm:"default:NULL"`
}

func GetAllVideos() ([]Video, error) {
	var videos []Video
	if err := common.DB.Preload("Files").Find(&videos).Error; err != nil {
		return nil, err
	}
	return videos, nil
}

func GetVideoByID(videoID uint) (*Video, error) {
	var video Video
	if err := common.DB.First(&video, videoID).Error; err != nil {
		return nil, err
	}
	return &video, nil
}
