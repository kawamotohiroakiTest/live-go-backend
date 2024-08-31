package handlers

import (
	"io"
	"live/common"
	"live/videoupload/models"
	"live/videoupload/services"
	"net/http"
	"os"
	"strconv"
)

func Upload(w http.ResponseWriter, r *http.Request) {

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		common.LogVideoUploadError(err)
		http.Error(w, "リクエストの解析に失敗しました", http.StatusBadRequest)
		return
	}

	// ユーザーIDの取得
	userID, err := common.GetUserIDFromContext(r.Context())
	if err != nil {
		common.LogVideoUploadError(err)
		http.Error(w, "ユーザー情報の取得に失敗しました", http.StatusUnauthorized)
		return
	}

	// DBトランザクションの開始
	tx := common.DB.Begin()
	if tx.Error != nil {
		common.LogVideoUploadError(tx.Error)
		http.Error(w, "内部サーバーエラー", http.StatusInternalServerError)
		return
	}

	// 動画情報の保存
	title := r.FormValue("title")
	description := r.FormValue("description")

	video, err := models.SaveVideoWithTransaction(tx, userID, title, description)
	if err != nil {
		common.LogVideoUploadError(err)
		tx.Rollback()
		http.Error(w, "動画情報の保存に失敗しました", http.StatusInternalServerError)
		return
	}

	// ENV_MODEの取得
	envMode := os.Getenv("ENV_MODE")

	var storageService *services.StorageService

	if envMode == "local" {
		storageService, err = services.InitMinioService()
	} else {
		storageService, err = services.NewStorageService()
	}

	if err != nil {
		common.LogVideoUploadError(err)
		tx.Rollback()
		http.Error(w, "ストレージサービスの初期化に失敗しました", http.StatusInternalServerError)
		return
	}

	// 動画ファイルの処理
	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		common.LogVideoUploadError(err)
		tx.Rollback()
		http.Error(w, "動画ファイルの取得に失敗しました", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// サムネイルファイルの処理
	thumbnail, thumbnailHeader, err := r.FormFile("thumbnail")
	if err != nil {
		common.LogVideoUploadError(err)
		tx.Rollback()
		http.Error(w, "サムネイルファイルの取得に失敗しました", http.StatusBadRequest)
		return
	}
	defer thumbnail.Close()

	// 動画ファイルのアップロード
	fileURL, err := storageService.UploadFile(file, fileHeader)
	if err != nil {
		common.LogVideoUploadError(err)
		tx.Rollback()
		http.Error(w, "動画ファイルのアップロードに失敗しました", http.StatusInternalServerError)
		return
	}

	// サムネイルファイルのアップロード
	thumbnailURL, err := storageService.UploadThumbnailFile(thumbnail, thumbnailHeader)
	if err != nil {
		common.LogVideoUploadError(err)
		tx.Rollback()
		http.Error(w, "サムネイルファイルのアップロードに失敗しました", http.StatusInternalServerError)
		return
	}

	// 動画ファイル情報の保存
	durationStr := r.FormValue("duration")
	duration, _ := strconv.Atoi(durationStr)
	fileSize := uint64(fileHeader.Size)
	format := fileHeader.Header.Get("Content-Type")

	_, err = models.SaveVideoFileWithTransaction(tx, video.ID, fileURL, thumbnailURL, uint(duration), fileSize, format)
	if err != nil {
		common.LogVideoUploadError(err)
		tx.Rollback()
		http.Error(w, "動画ファイル情報の保存に失敗しました", http.StatusInternalServerError)
		return
	}

	// トランザクションのコミット
	if err := tx.Commit().Error; err != nil {
		common.LogVideoUploadError(err)
		http.Error(w, "内部サーバーエラー", http.StatusInternalServerError)
		return
	}

	// 成功レスポンスを返す
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, `{"video_url":"`+fileURL+`", "thumbnail_url":"`+thumbnailURL+`"}`)
}
