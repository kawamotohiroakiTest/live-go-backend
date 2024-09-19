package comment

import (
	"fmt"
	"live/comment/handlers"
	"live/common"
	"net/http"

	"github.com/gorilla/mux"
)

func RegisterRoutes(router *mux.Router) {
	apiRouter := router.PathPrefix("/api/v1").Subrouter()

	// 認証が不要なルート
	// Get comments by video ID (認証不要)
	apiRouter.HandleFunc("/videos/{video_id}/comments", func(w http.ResponseWriter, r *http.Request) {
		handlers.GetCommentsByVideoIDHandler(w, r)
	}).Methods("GET")

	// 認証が必要なルートを設定
	authRouter := apiRouter.PathPrefix("").Subrouter() // 新しいサブルーター
	authRouter.Use(common.AuthMiddleware)              // AuthMiddleware を適用

	// Create a new comment (認証が必要)
	authRouter.HandleFunc("/comments", func(w http.ResponseWriter, r *http.Request) {
		handlers.CreateCommentHandler(w, r)
	}).Methods("POST")

	// Soft delete a comment by ID (認証が必要)
	authRouter.HandleFunc("/comments/{comment_id}", func(w http.ResponseWriter, r *http.Request) {
		handlers.DeleteCommentHandler(w, r)
	}).Methods("DELETE")

	fmt.Println("Registering comment routes")
}
