package main

import (
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"htnbk/internal/models"
)



func main() {
	if len(os.Args) < 4 {
		fmt.Println("Usage: htnbk <hatenaID> <blogID> <apiKey>")
		fmt.Println("Example: htnbk basyura blog.basyura.org your_api_key")
		os.Exit(1)
	}

	hatenaID := os.Args[1]
	blogID := os.Args[2]
	apiKey := os.Args[3]

	err := fetchAndSaveBlogEntries(hatenaID, blogID, apiKey)
	if err != nil {
		fmt.Printf("エラー: %v\n", err)
		os.Exit(1)
	}
}

func fetchAndSaveBlogEntries(hatenaID, blogID, apiKey string) error {
	allEntries, err := fetchAllBlogEntries(hatenaID, blogID, apiKey)
	if err != nil {
		return err
	}

	fmt.Printf("総記事数: %d\n", len(allEntries))
	fmt.Println(strings.Repeat("-", 50))

	for i, entry := range allEntries {
		publishedTime, err := time.Parse(time.RFC3339, entry.Published)
		if err != nil {
			fmt.Printf("日付解析エラー: %v\n", err)
			continue
		}
		fmt.Printf("%04d : %s - %s\n", i+1, publishedTime.Format("2006-01-02"), entry.Title)
		fmt.Printf("     ID: %s\n", entry.ID)

		// ファイルパスを生成（年/月/ファイル名）
		filePath, err := generateFilePath(entry.Published, entry.Title)
		if err != nil {
			fmt.Printf("ファイルパス生成エラー: %v\n", err)
			continue
		}

		// 必要なディレクトリを作成
		dir := filepath.Dir(filePath)
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			fmt.Printf("ディレクトリ作成エラー: %v\n", err)
			continue
		}

		// ファイルに保存
		err = saveEntryToFile(filePath, &entry)
		if err != nil {
			fmt.Printf("ファイル保存エラー: %v\n", err)
			continue
		}

		// 相対パスで表示
		relPath, _ := filepath.Rel(".", filePath)
		fmt.Printf("     保存: %s\n", relPath)
	}

	return nil
}

func fetchAllBlogEntries(hatenaID, blogID, apiKey string) ([]models.Entry, error) {
	var allEntries []models.Entry
	nextURL := fmt.Sprintf("https://blog.hatena.ne.jp/%s/%s/atom/entry", hatenaID, blogID)
	
	for nextURL != "" {
		fmt.Printf("取得中: %s\n", nextURL)
		
		entries, next, err := fetchBlogEntriesPage(hatenaID, apiKey, nextURL)
		if err != nil {
			return nil, err
		}
		
		allEntries = append(allEntries, entries...)
		nextURL = next
		
		fmt.Printf("  %d件取得（累計: %d件）\n", len(entries), len(allEntries))
	}
	
	return allEntries, nil
}

func fetchBlogEntriesPage(hatenaID, apiKey, url string) ([]models.Entry, string, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, "", err
	}

	// Basic認証のヘッダーを設定
	auth := base64.StdEncoding.EncodeToString([]byte(hatenaID + ":" + apiKey))
	req.Header.Set("Authorization", "Basic "+auth)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, "", fmt.Errorf("HTTPエラー: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", err
	}

	var feed models.Feed
	err = xml.Unmarshal(body, &feed)
	if err != nil {
		return nil, "", err
	}

	// nextリンクを検索
	var nextURL string
	for _, link := range feed.Links {
		if link.Rel == "next" {
			nextURL = link.Href
			break
		}
	}

	return feed.Entries, nextURL, nil
}


func generateFilePath(published, title string) (string, error) {
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

func saveEntryToFile(filePath string, entry *models.Entry) error {
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
