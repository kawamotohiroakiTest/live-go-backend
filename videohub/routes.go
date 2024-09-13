package videohub

import (
	"fmt"
	"live/common"
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
	videohubRouter.HandleFunc("/get_videos_by_ids", func(w http.ResponseWriter, r *http.Request) {
		handlers.GetVideosByIdsHandler(db, w, r)
	}).Methods("POST")
	fmt.Println("Registering videohub routes")

	videohubRouter.HandleFunc("/recommendations/user_{user_id}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		userID := vars["user_id"]

		fmt.Printf("Debug: Received user_id: %s\n", userID)

		common.LogVideoHubInfo(fmt.Sprintf("Debug: Received user_id: %s", userID))

		handlers.GetRecommendationsHandler(w, r)
	}).Methods("GET")

	fmt.Println("Registering videohub routes")
}
