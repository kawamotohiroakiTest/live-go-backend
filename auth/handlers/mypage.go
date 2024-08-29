package handlers

import (
	"encoding/json"
	"live/auth/models"
	"live/common"
	"net/http"
)

func MyPageHandler(w http.ResponseWriter, r *http.Request) {
	// クレームからユーザーIDを取得
	claims, ok := r.Context().Value("claims").(*common.Claims)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// ユーザー情報を取得
	var user models.User
	if err := common.DB.First(&user, claims.UserID).Error; err != nil {
		common.LogUser(common.ERROR, "User not found: "+err.Error())
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// ユーザーデータをJSON形式に変換
	userData, err := json.Marshal(user)
	if err != nil {
		common.LogUser(common.ERROR, "Failed to marshal user data: "+err.Error())
		http.Error(w, "Failed to retrieve user data", http.StatusInternalServerError)
		return
	}

	// user.logに書き込む
	common.LogUser(common.INFO, string(userData))

	// ユーザーデータをレスポンスとして返す
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(userData)
}
