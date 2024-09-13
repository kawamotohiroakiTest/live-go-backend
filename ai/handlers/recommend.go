package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/personalizeruntime"
)

// Recommend function for production environment
func Recommend(w http.ResponseWriter, r *http.Request) {
	// AWS設定を読み込み (IAMロールを使用)
	region := os.Getenv("PY_AWS_REGION")
	recommenderArn := os.Getenv("PY_RECOMMENDER_ARN")

	// ログに出力
	log.Printf("Region: %s, RecommenderArn: %s\n", region, recommenderArn)

	// AWS設定をロード (プロファイル指定はなし、IAMロールを使用)
	cfg, err := config.LoadDefaultConfig(r.Context(),
		config.WithRegion(region), // リージョンのみ指定
	)
	if err != nil {
		log.Printf("Error loading AWS config: %v", err)
		http.Error(w, "Failed to load AWS configuration", http.StatusInternalServerError)
		return
	}

	// Personalizeのクライアントを作成
	personalizeClient := personalizeruntime.NewFromConfig(cfg)

	// ユーザーIDをリクエストから取得
	userId := r.URL.Query().Get("user_id")
	if userId == "" {
		http.Error(w, "user_id is required", http.StatusBadRequest)
		return
	}

	// GetRecommendations APIを呼び出してレコメンデーションを取得
	input := &personalizeruntime.GetRecommendationsInput{
		RecommenderArn: aws.String(recommenderArn),
		UserId:         aws.String(userId),
		NumResults:     5,
	}

	result, err := personalizeClient.GetRecommendations(r.Context(), input)
	if err != nil {
		log.Printf("Error getting recommendations: %v", err)
		http.Error(w, fmt.Sprintf("failed to get recommendations: %v", err), http.StatusInternalServerError)
		return
	}

	// レコメンデーション結果をJSONで返す
	recommendations := make([]map[string]string, len(result.ItemList))
	for i, item := range result.ItemList {
		recommendations[i] = map[string]string{"itemId": *item.ItemId}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(recommendations)
}
