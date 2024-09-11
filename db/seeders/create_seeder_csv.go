package seeders

import (
	"encoding/csv"
	"fmt"
	"live/auth/models"
	videomodels "live/videohub/models"
	"os"
	"path/filepath"

	"gorm.io/gorm"
)

// Export users table to CSV
func ExportUsersToCSV(db *gorm.DB) error {
	var users []models.User
	if err := db.Find(&users).Error; err != nil {
		return fmt.Errorf("failed to fetch users: %v", err)
	}

	var records [][]string
	for _, user := range users {
		record := []string{
			fmt.Sprint(user.ID),       // USER_ID
			user.LastLoginAt.String(), // LAST_LOGIN
		}
		records = append(records, record)
	}

	headers := []string{"USER_ID", "LAST_LOGIN"}
	return createCSV("db/learningdata/users.csv", headers, records)
}

// Export videos table to CSV
func ExportVideosToCSV(db *gorm.DB) error {
	var videos []videomodels.Video
	if err := db.Find(&videos).Error; err != nil {
		return fmt.Errorf("failed to fetch videos: %v", err)
	}

	var records [][]string
	for _, video := range videos {
		record := []string{
			fmt.Sprint(video.ID),              // ITEM_ID
			video.Genre,                       // GENRES
			fmt.Sprint(video.Created.Unix()),  // CREATION_TIMESTAMP (UNIX timestamp)
			fmt.Sprint(video.ViewCount),       // VIEW_COUNT
			fmt.Sprintf("%.2f", video.Rating), // RATING
		}
		records = append(records, record)
	}

	headers := []string{"ITEM_ID", "GENRES", "CREATION_TIMESTAMP", "VIEW_COUNT", "RATING"}
	return createCSV("db/learningdata/videos.csv", headers, records)
}

// Export user_video_interactions table to CSV
func ExportUserVideoInteractionsToCSV(db *gorm.DB) error {
	var interactions []videomodels.UserVideoInteraction
	if err := db.Find(&interactions).Error; err != nil {
		return fmt.Errorf("failed to fetch user video interactions: %v", err)
	}

	var records [][]string
	for _, interaction := range interactions {
		record := []string{
			fmt.Sprint(interaction.UserID),           // USER_ID
			fmt.Sprint(interaction.VideoID),          // ITEM_ID
			interaction.EventType,                    // EVENT_TYPE
			fmt.Sprint(interaction.EventValue),       // EVENT_VALUE
			fmt.Sprint(interaction.CreatedAt.Unix()), // TIMESTAMP (UNIX timestamp)
		}
		records = append(records, record)
	}

	headers := []string{"USER_ID", "ITEM_ID", "EVENT_TYPE", "EVENT_VALUE", "TIMESTAMP"}
	return createCSV("db/learningdata/user_video_interactions.csv", headers, records)
}

// Helper function to create CSV file
func createCSV(filePath string, headers []string, records [][]string) error {
	dir := filepath.Dir(filePath)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
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

	// Write headers
	if err := writer.Write(headers); err != nil {
		return fmt.Errorf("failed to write CSV headers: %v", err)
	}

	// Write records
	for _, record := range records {
		if err := writer.Write(record); err != nil {
			return fmt.Errorf("failed to write record: %v", err)
		}
	}

	return nil
}
