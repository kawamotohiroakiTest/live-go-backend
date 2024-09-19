package handlers

import (
	"encoding/json"
	"fmt"
	"live/common"
	"live/videohub/models"
	"live/videohub/services"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func GetVideoByIDHandler(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	// URLからvideo_idを取得
	vars := mux.Vars(r)
	videoIDStr := vars["video_id"]

	// video_idをuintに変換
	videoID, err := strconv.ParseUint(videoIDStr, 10, 32)
	if err != nil {
		http.Error(w, "無効な動画IDです", http.StatusBadRequest)
		return
	}

	fmt.Printf("Debug: Fetching video with ID %d\n", videoID)

	// ストレージサービスの初期化
	var storageService *services.StorageService
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

	// 動画情報をデータベースから取得
	video, err := models.GetVideoByID(uint(videoID))
	if err != nil {
		http.Error(w, "動画が見つかりませんでした", http.StatusNotFound)
		return
	}

	// サムネイルや動画ファイルの署名付きURLを生成
	for i, file := range video.Files {
		if file.FilePath != "" {
			video.Files[i].FilePath, err = storageService.GetVideoPresignedURL(file.FilePath)
			if err != nil {
				common.LogVideoHubError(err)
				http.Error(w, "動画URLの取得に失敗しました", http.StatusInternalServerError)
				return
			}
		}
	}

	// 動画情報をJSONで返す
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(video); err != nil {
		common.LogVideoHubError(err)
		http.Error(w, "動画データのエンコードに失敗しました", http.StatusInternalServerError)
	}
}
