package common

import (
	"math/rand"
	"path/filepath"
	"time"

	"github.com/google/uuid"
)

// generateUniqueFileName は一意のファイル名を生成する関数です
func GenerateUniqueFileName(originalName string) string {
	extension := filepath.Ext(originalName)
	return uuid.New().String() + extension
}

// ランダムな日付を生成する共通関数
func RandomDate() time.Time {
	start := time.Date(time.Now().Year(), time.January, 1, 0, 0, 0, 0, time.Local)
	end := time.Now()
	delta := end.Unix() - start.Unix()

	sec := rand.Int63n(delta) + start.Unix()
	return time.Unix(sec, 0)
}

// ランダムな視聴回数を生成する関数
func RandomViewCount() int {
	return rand.Intn(100) + 1 // 1〜100の視聴回数をランダムに生成
}
