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

	// 動画ファイルの処理
	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		common.LogVideoUploadError(err)
		http.Error(w, "動画ファイルの取得に失敗しました", http.StatusBadRequest)
		return
	}
	defer file.Close()

	storageService, err := services.NewStorageService()
	if err != nil {
		common.LogVideoUploadError(err)
		http.Error(w, "ストレージサービスの初期化に失敗しました", http.StatusInternalServerError)
		return
	}

	// サムネイルファイルの処理
	thumbnail, thumbnailHeader, err := r.FormFile("thumbnail")
	if err != nil {
		common.LogVideoUploadError(err)
		http.Error(w, "サムネイルファイルの取得に失敗しました", http.StatusBadRequest)
		return
	}
	defer thumbnail.Close()

	// 動画ファイルのアップロード
	fileURL, err := storageService.UploadFile(file, fileHeader)
	if err != nil {
		common.LogVideoUploadError(err)
		http.Error(w, "動画ファイルのアップロードに失敗しました", http.StatusInternalServerError)
		return
	}

	// サムネイルファイルのアップロード
	thumbnailURL, err := storageService.UploadThumbnailFile(thumbnail, thumbnailHeader)
	if err != nil {
		common.LogVideoUploadError(err)
		http.Error(w, "サムネイルファイルのアップロードに失敗しました", http.StatusInternalServerError)
		return
	}

	// 成功レスポンスを返す
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, `{"video_url":"`+fileURL+`", "thumbnail_url":"`+thumbnailURL+`"}`)
}
