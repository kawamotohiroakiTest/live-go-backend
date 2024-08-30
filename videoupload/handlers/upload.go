package handlers

import (
	"io"
	"live/common"
	"live/videoupload/services"
	"net/http"
)

func Upload(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		common.LogVideoUploadError(err)
		http.Error(w, "リクエストの解析に失敗しました", http.StatusBadRequest)
		return
	}

	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		common.LogVideoUploadError(err)
		http.Error(w, "ファイルの取得に失敗しました", http.StatusBadRequest)
		return
	}
	defer file.Close()

	storageService, err := services.NewStorageService()
	if err != nil {
		common.LogVideoUploadError(err)
		http.Error(w, "ストレージサービスの初期化に失敗しました", http.StatusInternalServerError)
		return
	}

	fileURL, err := storageService.UploadFile(file, fileHeader)
	if err != nil {
		common.LogVideoUploadError(err)
		http.Error(w, "ファイルのアップロードに失敗しました", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, `{"url":"`+fileURL+`"}`)
}
