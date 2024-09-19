package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"live/comment/models"
	"live/common"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// CreateCommentHandler handles the creation of a new comment
func CreateCommentHandler(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		VideoID uint   `json:"video_id"`
		Content string `json:"content"`
	}

	// Log the incoming request for debugging
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}
	fmt.Printf("Request body: %s\n", string(body)) // Log raw request body

	// Now decode the request body into the payload struct
	if err := json.Unmarshal(body, &payload); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// ユーザーIDの取得
	userID, err := common.GetUserIDFromContext(r.Context())
	fmt.Println("userID: ", userID)
	if err != nil {
		common.LogVideoUploadError(err)
		http.Error(w, "ユーザー情報の取得に失敗しました", http.StatusUnauthorized)
		return
	}

	comment, err := models.CreateComment(userID, payload.VideoID, payload.Content)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Failed to create comment", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(comment)
}

// GetCommentsByVideoIDHandler retrieves comments for a specific video
func GetCommentsByVideoIDHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	videoIDStr := vars["video_id"]

	videoID, err := strconv.ParseUint(videoIDStr, 10, 32)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Invalid video ID", http.StatusBadRequest)
		return
	}

	comments, err := models.GetCommentsByVideoID(uint(videoID))
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Failed to fetch comments", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(comments)
}

// DeleteCommentHandler soft deletes a comment by setting the deleted timestamp
func DeleteCommentHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	commentIDStr := vars["comment_id"]

	commentID, err := strconv.ParseUint(commentIDStr, 10, 32)
	if err != nil {
		http.Error(w, "Invalid comment ID", http.StatusBadRequest)
		return
	}

	if err := models.DeleteComment(uint(commentID)); err != nil {
		http.Error(w, "Failed to delete comment", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	fmt.Fprintf(w, "Comment deleted successfully")
}
