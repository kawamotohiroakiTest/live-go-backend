package videohub

import (
	"live/videohub/handlers"
	"net/http"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func RegisterRoutes(router *mux.Router, db *gorm.DB) {
	apiRouter := router.PathPrefix("/api/v1").Subrouter()
	videohubRouter := apiRouter.PathPrefix("/videos").Subrouter()

	videohubRouter.HandleFunc("/list", handlers.ListVideos).Methods("GET")
	videohubRouter.HandleFunc("/create_user_video_interactions", func(w http.ResponseWriter, r *http.Request) {
		handlers.SaveUserVideoInteractionHandler(db, w, r)
	}).Methods("POST")
}
