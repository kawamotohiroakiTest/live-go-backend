package handlers

import (
	"encoding/json"
	"fmt"
	"live/common"
	"live/videohub/models"
	"live/videohub/services"
	"net/http"
	"os"
	"strings"

	"gorm.io/gorm"
)

// 動画検索
func SearchVideosHandler(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	// クエリパラメータから検索クエリを取得
	query := r.URL.Query().Get("query")
	if query == "" {
		http.Error(w, "検索クエリが必要です", http.StatusBadRequest)
		return
	}

	var videos []models.Video

	searchQuery := "%" + strings.ToLower(query) + "%"
	if err := db.Where("LOWER(title) LIKE ? OR LOWER(description) LIKE ?", searchQuery, searchQuery).Preload("Files").Find(&videos).Error; err != nil {
		http.Error(w, "動画の検索に失敗しました", http.StatusInternalServerError)
		return
	}

	// ストレージサービスの初期化
	var storageService *services.StorageService
	envMode := os.Getenv("ENV_MODE")
	var err error
	if envMode == "local" {
		storageService, err = services.InitMinioService()
	} else {
		storageService, err = services.NewStorageService()
	}
	if err != nil {
		common.LogVideoHubError(err)
		http.Error(w, "ストレージサービスの初期化に失敗しました", http.StatusInternalServerError)
		return
	}

	// 動画のファイルに署名付きURLを追加
	for i := range videos {
		for j, file := range videos[i].Files {
			if file.FilePath != "" {
				videos[i].Files[j].FilePath, err = storageService.GetVideoPresignedURL(file.FilePath)
				if err != nil {
					common.LogVideoHubError(err)
					http.Error(w, "動画URLの取得に失敗しました", http.StatusInternalServerError)
					return
				}
			}
		}
	}

	// デバッグ用に検索結果を出力
	fmt.Printf("検索クエリ: %s\n", query)
	fmt.Printf("検索結果: %+v\n", videos)

	// 結果をJSON形式で返す
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(videos); err != nil {
		http.Error(w, "結果のエンコードに失敗しました", http.StatusInternalServerError)
	}
}
