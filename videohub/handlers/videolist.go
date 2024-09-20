package handlers

import (
	"encoding/json"
	"fmt"
	"live/common"
	"live/videohub/models"
	"live/videohub/services"
	"net/http"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/personalizeruntime"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func ListVideos(w http.ResponseWriter, r *http.Request) {
	var storageService *services.StorageService
	var err error
	envMode := os.Getenv("ENV_MODE")
	if envMode == "local" {
		storageService, err = services.InitMinioService()
	} else {
		storageService, err = services.NewStorageService()
	}

	if err != nil {
		common.LogVideoHubError(err)
		http.Error(w, "ストレージサービスの初期化に失敗しました", http.StatusInternalServerError)
		return
	}

	videos, err := models.GetLimitedVideos(10)
	if err != nil {
		common.LogVideoHubError(err)
		http.Error(w, "動画の取得に失敗しました", http.StatusInternalServerError)
		return
	}

	// サムネイルと動画の署名付きURLを生成
	for _, video := range videos {
		for i, file := range video.Files {
			if file.FilePath != "" {
				// 署名付きURLを生成し、それを返す
				video.Files[i].FilePath, err = storageService.GetVideoPresignedURL(file.FilePath)
				if err != nil {
					common.LogVideoHubError(err)
					http.Error(w, "動画URLの取得に失敗しました", http.StatusInternalServerError)
					return
				}
			}
		}
	}

	// 動画データをJSONでフロントに返す
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(videos); err != nil {
		common.LogVideoHubError(err)
		http.Error(w, "動画データのエンコードに失敗しました", http.StatusInternalServerError)
	}
}

func GetVideosByIdsHandler(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	// ストレージサービスの初期化
	var storageService *services.StorageService
	var err error
	envMode := os.Getenv("ENV_MODE")
	if envMode == "local" {
		storageService, err = services.InitMinioService()
	} else {
		storageService, err = services.NewStorageService()
	}

	if err != nil {
		common.LogVideoHubError(err)
		http.Error(w, "ストレージサービスの初期化に失敗しました", http.StatusInternalServerError)
		return
	}

	// リクエストボディから動画IDリストを受け取る
	var requestBody struct {
		VideoIds []string `json:"videoIds"`
	}

	// リクエストボディのデコード
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if len(requestBody.VideoIds) == 0 {
		http.Error(w, "No video IDs provided", http.StatusBadRequest)
		return
	}

	// "video_1" 形式からIDだけを抽出
	var ids []uint
	for _, videoId := range requestBody.VideoIds {
		idStr := strings.TrimPrefix(videoId, "video_")
		idStr = strings.TrimPrefix(idStr, "item_")

		var id uint
		fmt.Sscanf(idStr, "%d", &id) // 文字列を数値に変換
		ids = append(ids, id)
	}
	fmt.Println("ids", ids)

	// 複数の動画IDに基づいて動画情報を取得するメソッドを呼び出す
	videos, err := models.GetVideosByIds(ids)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error retrieving videos: %v", err), http.StatusInternalServerError)
		return
	}

	// サムネイルと動画の署名付きURLを生成
	for _, video := range videos {
		for i, file := range video.Files {
			// 動画URLの生成
			if file.FilePath != "" {
				video.Files[i].FilePath, err = storageService.GetVideoPresignedURL(file.FilePath)
				fmt.Println("video.Files[i].FilePathvideolistAI", video.Files[i].FilePath)
				if err != nil {
					common.LogVideoHubError(err)
					http.Error(w, "動画URLの取得に失敗しました", http.StatusInternalServerError)
					return
				}
			}
		}
	}

	// 取得した動画情報をJSONで返す
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(videos); err != nil {
		common.LogVideoHubError(err)
		http.Error(w, "動画のJSON変換に失敗しました", http.StatusInternalServerError)
		return
	}
}

func GetRecommendationsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("GetRecommendationsHandler")
	// リクエストパラメータからuser_idを取得
	vars := mux.Vars(r)
	userID := vars["user_id"]
	fmt.Printf("Received user_id: %s\n", userID)

	// 環境変数からAWSの設定を取得
	region := os.Getenv("PY_AWS_REGION")
	recommenderArn := os.Getenv("PY_RECOMMENDER_ARN")

	if region == "" || recommenderArn == "" {
		fmt.Println("AWSの設定が正しくありません")
		http.Error(w, "AWSの設定が正しくありません", http.StatusInternalServerError)
		return
	}

	// AWSセッションの作成
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})
	if err != nil {
		fmt.Println("AWSセッションの初期化に失敗しました: %v", err)
		common.LogVideoHubError(fmt.Errorf("AWSセッションの初期化に失敗しました: %v", err))
		http.Error(w, "AWSセッションの初期化に失敗しました", http.StatusInternalServerError)
		return
	}

	// Personalize Runtimeのクライアントを作成
	personalizeRuntime := personalizeruntime.New(sess)

	// GetRecommendations APIを呼び出してレコメンデーションを取得
	input := &personalizeruntime.GetRecommendationsInput{
		RecommenderArn: aws.String(recommenderArn),
		UserId:         aws.String(userID),
		NumResults:     aws.Int64(5), // 推奨結果を5件に制限
	}

	result, err := personalizeRuntime.GetRecommendations(input)
	if err != nil {
		fmt.Println("レコメンデーション取得に失敗しました: %v", err)
		common.LogVideoHubError(fmt.Errorf("レコメンデーション取得に失敗しました: %v", err))
		http.Error(w, "レコメンデーション取得に失敗しました", http.StatusInternalServerError)
		return
	}

	// レコメンデーション結果をJSONで返す
	recommendations := make([]map[string]string, len(result.ItemList))
	for i, item := range result.ItemList {
		recommendations[i] = map[string]string{"itemId": *item.ItemId}
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(recommendations); err != nil {
		common.LogVideoHubError(fmt.Errorf("レコメンデーションのJSONエンコードに失敗しました: %v", err))
		http.Error(w, "レコメンデーションのJSONエンコードに失敗しました", http.StatusInternalServerError)
		return
	}
}

func StreamVideoHandler(w http.ResponseWriter, r *http.Request) {
	// ストレージサービスの初期化
	var storageService *services.StorageService
	var err error
	envMode := os.Getenv("ENV_MODE")
	if envMode == "local" {
		storageService, err = services.InitMinioService()
	} else {
		storageService, err = services.NewStorageService()
	}

	if err != nil {
		common.LogVideoHubError(err)
		http.Error(w, "ストレージサービスの初期化に失敗しました", http.StatusInternalServerError)
		return
	}

	// フロントエンドから送信された動画ファイル名を取得
	vars := mux.Vars(r)
	filePath := vars["filename"]

	// 動画ファイルの署名付きURLを取得
	videoURL, err := storageService.GetVideoPresignedURL(filePath)
	if err != nil {
		http.Error(w, "署名付きURLの取得に失敗しました", http.StatusInternalServerError)
		return
	}

	// 動画をクライアントにストリーミング送信
	http.Redirect(w, r, videoURL, http.StatusTemporaryRedirect) // 署名付きURLを返してリダイレクト
}
