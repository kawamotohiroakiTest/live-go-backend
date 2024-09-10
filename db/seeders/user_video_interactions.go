package seeders

import (
	"fmt"
	"live/auth/models"
	videomodels "live/videohub/models"
	"log"
	"math/rand"
	"time"

	"gorm.io/gorm"
)

var eventTypes = []string{"play", "pause", "complete", "like", "dislike"}

func SeedUserVideoInteractions(db *gorm.DB) {
	var users []models.User
	var videos []videomodels.Video

	// 全てのユーザーとビデオを取得
	if err := db.Find(&users).Error; err != nil {
		log.Printf("Failed to retrieve users: %v", err)
		return
	}

	if err := db.Find(&videos).Error; err != nil {
		log.Printf("Failed to retrieve videos: %v", err)
		return
	}

	for i := 1; i <= 100; i++ {
		// ランダムなユーザーとビデオを選択
		user := users[rand.Intn(len(users))]
		video := videos[rand.Intn(len(videos))]
		eventType := eventTypes[rand.Intn(len(eventTypes))]

		interaction := videomodels.UserVideoInteraction{
			UserID:    user.ID,
			VideoID:   video.ID,
			EventType: eventType,
			CreatedAt: time.Now(),
		}

		if err := db.Create(&interaction).Error; err != nil {
			log.Printf("Failed to create user video interaction: %v", err)
		} else {
			fmt.Printf("Created user video interaction: UserID=%d, VideoID=%d, EventType=%s\n", user.ID, video.ID, eventType)
		}
	}
}
