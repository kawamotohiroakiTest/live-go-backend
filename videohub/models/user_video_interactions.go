package models

import "time"

type UserVideoInteraction struct {
	ID         uint      `gorm:"primaryKey"`
	UserID     uint      `json:"user_id" gorm:"not null"`                                                              // JSONタグを追加
	VideoID    uint      `json:"video_id" gorm:"not null"`                                                             // JSONタグを追加
	EventType  string    `json:"event_type" gorm:"type:enum('play', 'pause', 'complete', 'like', 'dislike');not null"` // JSONタグを追加
	EventValue float64   `json:"event_value" gorm:"nullable"`                                                          // オプション: 視聴時間や評価スコアなど
	CreatedAt  time.Time `gorm:"autoCreateTime"`                                                                       // 作成日時
}
