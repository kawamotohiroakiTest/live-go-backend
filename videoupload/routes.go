package videoupload

import (
	"live/common"
	"live/videoupload/handlers"

	"github.com/gorilla/mux"
)

func RegisterRoutes(router *mux.Router) {
	videouploadRouter := router.PathPrefix("/api/v1/videoupload").Subrouter()
	videouploadRouter.Use(common.AuthMiddleware)

	videouploadRouter.HandleFunc("/upload", handlers.Upload).Methods("POST")
}
