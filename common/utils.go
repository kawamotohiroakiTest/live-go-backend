package common

import (
	"path/filepath"

	"github.com/google/uuid"
)

// generateUniqueFileName は一意のファイル名を生成する関数です
func GenerateUniqueFileName(originalName string) string {
	extension := filepath.Ext(originalName)
	return uuid.New().String() + extension
}
