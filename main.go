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
)

type Feed struct {
	XMLName xml.Name `xml:"feed"`
	Title   string   `xml:"title"`
	Links   []Link   `xml:"link"`
	Entries []Entry  `xml:"entry"`
}

type Link struct {
	Rel  string `xml:"rel,attr"`
	Href string `xml:"href,attr"`
}

type Entry struct {
	Title     string `xml:"title"`
	Published string `xml:"published"`
	Updated   string `xml:"updated"`
	Links     []Link `xml:"link"`
	Content   string `xml:"content"`
}

type DetailedEntry struct {
	XMLName xml.Name `xml:"entry"`
	Title   string   `xml:"title"`
	Content string   `xml:"content"`
	Published string `xml:"published"`
	Updated   string `xml:"updated"`
}

func main() {
	apiKey := os.Getenv("HATENA_API_KEY")
	if apiKey == "" {
		fmt.Println("HATENA_API_KEY環境変数が設定されていません")
		os.Exit(1)
	}

	// はてなIDとブログIDを設定（後で引数から取得するように変更可能）
	hatenaID := "basyura"
	blogID := "blog.basyura.org"

	err := fetchAndSaveBlogEntries(hatenaID, blogID, apiKey)
	if err != nil {
		fmt.Printf("エラー: %v\n", err)
		os.Exit(1)
	}
}

func fetchAndSaveBlogEntries(hatenaID, blogID, apiKey string) error {
	entries, err := fetchBlogEntries(hatenaID, blogID, apiKey)
	if err != nil {
		return err
	}

	// 出力ディレクトリを作成
	outputDir := "entries"
	err = os.MkdirAll(outputDir, 0755)
	if err != nil {
		return fmt.Errorf("ディレクトリ作成エラー: %v", err)
	}

	for i, entry := range entries {
		publishedTime, err := time.Parse(time.RFC3339, entry.Published)
		if err != nil {
			fmt.Printf("日付解析エラー: %v\n", err)
			continue
		}
		fmt.Printf("%04d : %s - %s\n", i+1, publishedTime.Format("2006-01-02"), entry.Title)

		// 個別エントリーの詳細を取得
		detailedEntry, err := fetchEntryDetail(hatenaID, blogID, apiKey, entry)
		if err != nil {
			fmt.Printf("エントリー取得エラー: %v\n", err)
			continue
		}

		// ファイル名を生成
		fileName := generateFileName(i+1, entry.Published, entry.Title)
		filePath := filepath.Join(outputDir, fileName)

		// ファイルに保存
		err = saveEntryToFile(filePath, detailedEntry)
		if err != nil {
			fmt.Printf("ファイル保存エラー: %v\n", err)
			continue
		}

		fmt.Printf("     保存: %s\n", fileName)
	}

	return nil
}

func fetchBlogEntries(hatenaID, blogID, apiKey string) ([]Entry, error) {
	url := fmt.Sprintf("https://blog.hatena.ne.jp/%s/%s/atom/entry", hatenaID, blogID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// Basic認証のヘッダーを設定
	auth := base64.StdEncoding.EncodeToString([]byte(hatenaID + ":" + apiKey))
	req.Header.Set("Authorization", "Basic "+auth)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTPエラー: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var feed Feed
	err = xml.Unmarshal(body, &feed)
	if err != nil {
		return nil, err
	}

	return feed.Entries, nil
}

func fetchEntryDetail(hatenaID, blogID, apiKey string, entry Entry) (*DetailedEntry, error) {
	var editLink string
	for _, link := range entry.Links {
		if link.Rel == "edit" {
			editLink = link.Href
			break
		}
	}

	if editLink == "" {
		return nil, fmt.Errorf("edit linkが見つかりません")
	}

	req, err := http.NewRequest("GET", editLink, nil)
	if err != nil {
		return nil, err
	}

	auth := base64.StdEncoding.EncodeToString([]byte(hatenaID + ":" + apiKey))
	req.Header.Set("Authorization", "Basic "+auth)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTPエラー: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var detailedEntry DetailedEntry
	err = xml.Unmarshal(body, &detailedEntry)
	if err != nil {
		return nil, err
	}

	return &detailedEntry, nil
}

func generateFileName(index int, published, title string) string {
	publishedTime, err := time.Parse(time.RFC3339, published)
	if err != nil {
		publishedTime = time.Now()
	}

	// タイトルをファイル名に適したフォーマットに変換
	safeTitle := strings.ReplaceAll(title, " ", "_")
	// ファイルシステムで禁止されている文字を置き換え
	invalidChars := []string{"/", "\\", ":", "*", "?", "\"", "<", ">", "|"}
	for _, char := range invalidChars {
		safeTitle = strings.ReplaceAll(safeTitle, char, "_")
	}
	// ファイル名の長さ制限（バイト数で制限）
	for len(safeTitle) > 50 {
		// 最後の文字を削除
		runes := []rune(safeTitle)
		if len(runes) > 0 {
			safeTitle = string(runes[:len(runes)-1])
		} else {
			break
		}
	}

	return fmt.Sprintf("%s_%s.md", publishedTime.Format("2006-01-02"), safeTitle)
}

func saveEntryToFile(filePath string, entry *DetailedEntry) error {
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
