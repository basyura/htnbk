package storage

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"htnblg-export/internal/models"
)

// GenerateFilePath は記事の公開日とタイトルから適切なファイルパスを生成する
func GenerateFilePath(published, title string) (string, error) {
	publishedTime, err := time.Parse(time.RFC3339, published)
	if err != nil {
		return "", err
	}

	// タイトルをファイル名に適したフォーマットに変換
	safeTitle := strings.ReplaceAll(title, " ", "_")
	// ファイルシステムで禁止されている文字を置き換え
	invalidChars := []string{"/", "\\", ":", "*", "?", "\"", "<", ">", "|"}
	for _, char := range invalidChars {
		safeTitle = strings.ReplaceAll(safeTitle, char, "_")
	}

	// 年/月/日付_タイトル.md の形式
	year := publishedTime.Format("2006")
	month := publishedTime.Format("01")
	fileName := fmt.Sprintf("%s_%s.md", publishedTime.Format("2006-01-02"), safeTitle)

	return filepath.Join("entries", year, month, fileName), nil
}

// SaveEntryToFile は記事をMarkdownファイルとして保存する
func SaveEntryToFile(filePath string, entry *models.Entry) error {
	publishedTime, err := time.Parse(time.RFC3339, entry.Published)
	if err != nil {
		publishedTime = time.Now()
	}

	content := fmt.Sprintf(`---
title: %s
published: %s
updated: %s
---

%s
`, entry.Title, publishedTime.Format("2006-01-02"), entry.Updated, entry.Content)

	return os.WriteFile(filePath, []byte(content), 0644)
}

// GetLatestEntryDate はentriesフォルダから最新の記事日付を取得する
func GetLatestEntryDate() (time.Time, error) {
	var latestDate time.Time

	entriesDir := "entries"
	if _, err := os.Stat(entriesDir); os.IsNotExist(err) {
		// entriesフォルダが存在しない場合は、最古の日付を返す
		return time.Time{}, nil
	}

	err := filepath.Walk(entriesDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(info.Name(), ".md") {
			// ファイル名から日付を抽出 (YYYY-MM-DD_title.md形式)
			fileName := info.Name()
			if len(fileName) >= 10 {
				dateStr := fileName[:10] // YYYY-MM-DD部分
				if fileDate, err := time.Parse("2006-01-02", dateStr); err == nil {
					if fileDate.After(latestDate) {
						latestDate = fileDate
					}
				}
			}
		}
		return nil
	})

	return latestDate, err
}
