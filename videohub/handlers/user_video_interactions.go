package handlers

import (
	"encoding/json"
	"fmt"
	"live/common"
	"live/videohub/models"
	"live/videohub/services"
	"net/http"

	"gorm.io/gorm"
)

// SaveUserVideoInteractionHandler ハンドラー
func SaveUserVideoInteractionHandler(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	var interaction models.UserVideoInteraction

	// リクエストボディをJSONとしてデコード
	if err := json.NewDecoder(r.Body).Decode(&interaction); err != nil {
		common.LogVideoHubError(err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	common.LogVideoHubError(fmt.Errorf("Received event_type: %s", interaction.EventType))

	// サービスを使ってインタラクションを保存
	if err := services.SaveUserVideoInteraction(db, &interaction); err != nil {
		common.LogVideoHubError(err)
		http.Error(w, fmt.Sprintf("Failed to save interaction: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Interaction saved successfully"})
}
