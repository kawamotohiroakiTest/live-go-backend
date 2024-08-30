package models

import (
	"live/common"
	"time"
)

type Video struct {
	ID          uint       `gorm:"primary_key"`
	UserID      uint       `gorm:"not null"`
	Title       string     `gorm:"type:varchar(255);not null"`
	Description string     `gorm:"type:text"`
	Created     time.Time  `gorm:"default:CURRENT_TIMESTAMP"`
	Modified    time.Time  `gorm:"default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
	Deleted     *time.Time `gorm:"default:NULL"`
}

type VideoFile struct {
	ID            uint       `gorm:"primary_key"`
	VideoID       uint       `gorm:"not null"`
	FilePath      string     `gorm:"type:varchar(255);not null"`
	ThumbnailPath string     `gorm:"type:varchar(255)"`
	Duration      uint       `gorm:"type:int"`
	FileSize      uint64     `gorm:"type:bigint"`
	Format        string     `gorm:"type:varchar(50);not null"`
	Status        string     `gorm:"type:enum('pending','processing','completed','failed');default:'pending'"`
	Created       time.Time  `gorm:"default:CURRENT_TIMESTAMP"`
	Modified      time.Time  `gorm:"default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
	Deleted       *time.Time `gorm:"default:NULL"`
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
