package storage

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
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

	// エントリーのURLを取得
	var entryURL string
	for _, link := range entry.Links {
		if link.Rel == "alternate" {
			entryURL = link.Href
			break
		}
	}

	// タグとカテゴリを取得
	var tags []string
	for _, category := range entry.Categories {
		if category.Term != "" {
			tags = append(tags, category.Term)
		}
	}

	// ヘッダーを構築
	content := fmt.Sprintf(`---
title: %s
published: %s
updated: %s
url: %s
entry-id: %s`, entry.Title, publishedTime.Format("2006-01-02"), entry.Updated, entryURL, entry.ID)

	// 作成者情報
	if entry.Author.Name != "" {
		content += fmt.Sprintf("\nauthor: %s", entry.Author.Name)
	}

	// 編集日時
	if entry.Edited != "" {
		content += fmt.Sprintf("\nedited: %s", entry.Edited)
	}

	// 下書き状態
	if entry.Control.Draft == "yes" {
		content += "\ndraft: true"
	}

	// プレビュー共有状態
	if entry.Control.Preview == "yes" {
		content += "\npreview: true"
	}

	// カスタムURL
	if entry.CustomURL != "" {
		content += fmt.Sprintf("\ncustom-url: %s", entry.CustomURL)
	}

	// タグがある場合は追加
	if len(tags) > 0 {
		content += "\ntags:\n"
		for _, tag := range tags {
			content += fmt.Sprintf("  - %s\n", tag)
		}
	}

	content += "---\n\n" + entry.Content + "\n"

	return os.WriteFile(filePath, []byte(content), 0644)
}

// GetLatestEntryDate はentriesフォルダから最新の記事日付を取得する
func GetLatestEntryDate() (time.Time, error) {
	entriesDir := "entries"
	if _, err := os.Stat(entriesDir); os.IsNotExist(err) {
		// entriesフォルダが存在しない場合は、最古の日付を返す
		return time.Time{}, nil
	}

	// 年ディレクトリを取得して降順ソート
	yearDirs, err := os.ReadDir(entriesDir)
	if err != nil {
		return time.Time{}, err
	}

	// 年ディレクトリ名を収集
	var years []string
	for _, yearDir := range yearDirs {
		if yearDir.IsDir() {
			years = append(years, yearDir.Name())
		}
	}
	if len(years) == 0 {
		return time.Time{}, nil
	}

	// 降順ソート（新しい年から）
	sort.Sort(sort.Reverse(sort.StringSlice(years)))

	// 最新の年から順に検索
	for _, year := range years {
		latestInYear, err := getLatestDateInYear(filepath.Join(entriesDir, year))
		if err != nil {
			continue // エラーがあっても次の年を試す
		}
		if !latestInYear.IsZero() {
			return latestInYear, nil // 最初に見つかった年の最新日付を返す
		}
	}

	return time.Time{}, nil
}

// getLatestDateInYear は指定された年ディレクトリ内の最新日付を取得
func getLatestDateInYear(yearDir string) (time.Time, error) {
	// 月ディレクトリを取得
	monthDirs, err := os.ReadDir(yearDir)
	if err != nil {
		return time.Time{}, err
	}

	// 月ディレクトリ名を収集
	var months []string
	for _, monthDir := range monthDirs {
		if monthDir.IsDir() {
			months = append(months, monthDir.Name())
		}
	}
	if len(months) == 0 {
		return time.Time{}, nil
	}

	// 降順ソート（新しい月から）
	sort.Sort(sort.Reverse(sort.StringSlice(months)))

	// 最新の月から順に検索
	for _, month := range months {
		latestInMonth, err := getLatestDateInMonth(filepath.Join(yearDir, month))
		if err != nil {
			continue // エラーがあっても次の月を試す
		}
		if !latestInMonth.IsZero() {
			return latestInMonth, nil // 最初に見つかった月の最新日付を返す
		}
	}

	return time.Time{}, nil
}

// getLatestDateInMonth は指定された月ディレクトリ内の最新日付を取得
func getLatestDateInMonth(monthDir string) (time.Time, error) {
	files, err := os.ReadDir(monthDir)
	if err != nil {
		return time.Time{}, err
	}

	// ファイル名を収集（日付順でソートするため）
	var fileNames []string
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".md") && len(file.Name()) >= 10 {
			fileNames = append(fileNames, file.Name())
		}
	}
	if len(fileNames) == 0 {
		return time.Time{}, nil
	}

	// ファイル名を降順ソート（YYYY-MM-DDの部分で新しい日付から）
	sort.Sort(sort.Reverse(sort.StringSlice(fileNames)))

	// 最初のファイル（最新日付）の日付を解析して返す
	fileName := fileNames[0]
	dateStr := fileName[:10] // YYYY-MM-DD部分
	fileDate, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return time.Time{}, err
	}

	return fileDate, nil
}
