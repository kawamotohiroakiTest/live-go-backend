package services

import (
	"live/common"
	"live/videohub/models"

	"gorm.io/gorm"
)

// user_video_interactionsに新しいレコードを保存
func SaveUserVideoInteraction(db *gorm.DB, interaction *models.UserVideoInteraction) error {
	if err := db.Create(interaction).Error; err != nil {
		common.LogVideoHubError(err)
		return err
	}
	return nil
}
