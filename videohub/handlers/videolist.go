package handlers

import (
	"encoding/json"
	"live/common"
	"live/videohub/models"
	"live/videohub/services"
	"net/http"
	"os"
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
			// サムネイルURLの生成
			// if file.ThumbnailPath != "" {
			// 	video.Files[i].ThumbnailPath, err = storageService.GetThumbnailPresignedURL(file.ThumbnailPath)
			// 	if err != nil {
			// 		common.LogVideoHubError(err)
			// 		http.Error(w, "サムネイルURLの取得に失敗しました", http.StatusInternalServerError)
			// 		return
			// 	}
			// }

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
