package handlers

import (
	"encoding/json"
	"live/common"
	"live/videohub/models"
	"net/http"
)

func ListVideos(w http.ResponseWriter, r *http.Request) {
	videos, err := models.GetAllVideos()
	if err != nil {
		common.LogVideoHubError(err)
		http.Error(w, "動画の取得に失敗しました", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(videos); err != nil {
		common.LogVideoHubError(err)
		http.Error(w, "動画のJSON変換に失敗しました", http.StatusInternalServerError)
		return
	}
}
