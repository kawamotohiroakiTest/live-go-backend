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

func ListVideos(w http.ResponseWriter, r *http.Request) {
	// ストレージサービスの初期化
	var storageService *services.StorageService
	var err error
	envMode := os.Getenv("ENV_MODE")
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

	videos, err := models.GetAllVideos()
	if err != nil {
		common.LogVideoHubError(err)
		http.Error(w, "動画の取得に失敗しました", http.StatusInternalServerError)
		return
	}

	// サムネイルと動画の署名付きURLを生成
	for _, video := range videos {
		for i, file := range video.Files {
			// 動画URLの生成
			if file.FilePath != "" {
				video.Files[i].FilePath, err = storageService.GetVideoPresignedURL(file.FilePath)

				if err != nil {
					common.LogVideoHubError(err)
					http.Error(w, "動画URLの取得に失敗しました", http.StatusInternalServerError)
					return
				}
			}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(videos); err != nil {
		common.LogVideoHubError(err)
		http.Error(w, "動画のJSON変換に失敗しました", http.StatusInternalServerError)
		return
	}
}

func GetVideosByIdsHandler(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	// ストレージサービスの初期化
	var storageService *services.StorageService
	var err error
	envMode := os.Getenv("ENV_MODE")
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

	// リクエストボディから動画IDリストを受け取る
	var requestBody struct {
		VideoIds []string `json:"videoIds"`
	}

	// リクエストボディのデコード
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if len(requestBody.VideoIds) == 0 {
		http.Error(w, "No video IDs provided", http.StatusBadRequest)
		return
	}

	// "video_1" 形式からIDだけを抽出
	var ids []uint
	for _, videoId := range requestBody.VideoIds {
		idStr := strings.TrimPrefix(videoId, "video_")
		var id uint
		fmt.Sscanf(idStr, "%d", &id) // 文字列を数値に変換
		ids = append(ids, id)
	}
	fmt.Println("ids", ids)

	// 複数の動画IDに基づいて動画情報を取得するメソッドを呼び出す
	videos, err := models.GetVideosByIds(ids)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error retrieving videos: %v", err), http.StatusInternalServerError)
		return
	}

	// サムネイルと動画の署名付きURLを生成
	for _, video := range videos {
		for i, file := range video.Files {
			// 動画URLの生成
			if file.FilePath != "" {
				video.Files[i].FilePath, err = storageService.GetVideoPresignedURL(file.FilePath)
				fmt.Println("video.Files[i].FilePathvideolistAI", video.Files[i].FilePath)
				if err != nil {
					common.LogVideoHubError(err)
					http.Error(w, "動画URLの取得に失敗しました", http.StatusInternalServerError)
					return
				}
			}
		}
	}

	// 取得した動画情報をJSONで返す
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(videos); err != nil {
		common.LogVideoHubError(err)
		http.Error(w, "動画のJSON変換に失敗しました", http.StatusInternalServerError)
		return
	}
}
