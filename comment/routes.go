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

	// Create a new comment
	// 認証が必要なルート
	apiRouter.Use(common.AuthMiddleware) // AuthMiddleware を適用

	// Create a new comment (認証が必要)
	apiRouter.HandleFunc("/comments", func(w http.ResponseWriter, r *http.Request) {
		handlers.CreateCommentHandler(w, r)
	}).Methods("POST")

	// Get comments by video ID
	apiRouter.HandleFunc("/videos/{video_id}/comments", func(w http.ResponseWriter, r *http.Request) {
		handlers.GetCommentsByVideoIDHandler(w, r)
	}).Methods("GET")

	// Soft delete a comment by ID
	apiRouter.HandleFunc("/comments/{comment_id}", func(w http.ResponseWriter, r *http.Request) {
		handlers.DeleteCommentHandler(w, r)
	}).Methods("DELETE")

	fmt.Println("Registering comment routes")
}
