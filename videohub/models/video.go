package models

import (
	"live/common"
	"time"
)

type Video struct {
	ID            uint       `gorm:"primary_key"`
	Title         string     `gorm:"type:varchar(255);not null"`
	Description   string     `gorm:"type:text"`
	ThumbnailPath string     `gorm:"type:varchar(255)"`
	VideoPath     string     `gorm:"type:varchar(255);not null"`
	Created       time.Time  `gorm:"default:CURRENT_TIMESTAMP"`
	Modified      time.Time  `gorm:"default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
	Deleted       *time.Time `gorm:"default:NULL"`
}

func GetAllVideos() ([]Video, error) {
	var videos []Video
	if err := common.DB.Find(&videos).Error; err != nil {
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
