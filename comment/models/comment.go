package models

import (
	"live/common"
	"time"
)

type Comment struct {
	ID       uint       `gorm:"primary_key"`
	UserID   uint       `gorm:"not null"`                                              // 外部キー usersテーブルのID
	VideoID  uint       `gorm:"not null"`                                              // 外部キー videosテーブルのID
	Content  string     `gorm:"type:varchar(255);not null"`                            // コメント内容
	Created  time.Time  `gorm:"default:CURRENT_TIMESTAMP"`                             // 作成日時
	Modified time.Time  `gorm:"default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"` // 更新日時
	Deleted  *time.Time `gorm:"default:NULL"`                                          // 削除日時
}

// CreateComment creates a new comment in the database
func CreateComment(userID, videoID uint, content string) (*Comment, error) {
	comment := &Comment{
		UserID:  userID,
		VideoID: videoID,
		Content: content,
	}
	if err := common.DB.Create(comment).Error; err != nil {
		return nil, err
	}
	return comment, nil
}

// GetCommentsByVideoID retrieves comments for a specific video
func GetCommentsByVideoID(videoID uint) ([]Comment, error) {
	var comments []Comment
	if err := common.DB.Where("video_id = ? AND deleted IS NULL", videoID).Find(&comments).Error; err != nil {
		return nil, err
	}
	return comments, nil
}

// DeleteComment soft deletes a comment by setting the deleted timestamp
func DeleteComment(commentID uint) error {
	var comment Comment
	if err := common.DB.First(&comment, commentID).Error; err != nil {
		return err
	}
	now := time.Now()
	comment.Deleted = &now
	return common.DB.Save(&comment).Error
}
