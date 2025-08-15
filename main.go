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

	"htnblg-export/internal/models"
)

func main() {
	if len(os.Args) < 4 {
		fmt.Println("Usage: htnblg-export [--all] <hatenaID> <blogID> <apiKey>")
		fmt.Println("       --all: すべての記事を取得（デフォルトは最新日付以降のみ）")
		fmt.Println("Example: htnblg-export basyura blog.basyura.org your_api_key")
		fmt.Println("         htnblg-export --all basyura blog.basyura.org your_api_key")
		os.Exit(1)
	}

	var allMode bool
	var hatenaID, blogID, apiKey string

	// --allオプションをチェック
	if os.Args[1] == "--all" {
		if len(os.Args) < 5 {
			fmt.Println("Usage: htnblg-export --all <hatenaID> <blogID> <apiKey>")
			os.Exit(1)
		}
		allMode = true
		hatenaID = os.Args[2]
		blogID = os.Args[3]
		apiKey = os.Args[4]
	} else {
		hatenaID = os.Args[1]
		blogID = os.Args[2]
		apiKey = os.Args[3]
	}

	err := doMain(hatenaID, blogID, apiKey, allMode)
	if err != nil {
		fmt.Printf("エラー: %v\n", err)
		os.Exit(1)
	}
}

func doMain(hatenaID, blogID, apiKey string, allMode bool) error {
	// 最新の記事日付を取得（--allモードでない場合）
	var sinceDate time.Time
	if !allMode {
		var err error
		sinceDate, err = getLatestEntryDate()
		if err != nil {
			fmt.Printf("警告: 最新記事日付の取得に失敗: %v\n", err)
		} else if !sinceDate.IsZero() {
			fmt.Printf("最新記事日付: %s\n", sinceDate.Format("2006-01-02"))
			fmt.Printf("この日付以降の記事のみを取得します。すべて取得する場合は --all オプションを使用してください。\n")
		} else {
			fmt.Printf("初回実行のため、すべての記事を取得します。\n")
		}
	} else {
		fmt.Printf("--all モードですべての記事を取得します。\n")
	}

	// 効率的な取得：APIレベルで必要な分だけ取得
	allEntries, err := fetchAllBlogEntries(hatenaID, blogID, apiKey, allMode, sinceDate)
	if err != nil {
		return err
	}

	// 取得結果の表示
	if !allMode && !sinceDate.IsZero() {
		fmt.Printf("新しい記事: %d件\n", len(allEntries))
	} else {
		fmt.Printf("総記事数: %d件\n", len(allEntries))
	}

	if len(allEntries) == 0 {
		fmt.Println("新しい記事はありません。")
		return nil
	}

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

func fetchAllBlogEntries(hatenaID, blogID, apiKey string, allMode bool, sinceDate time.Time) ([]models.Entry, error) {
	var allEntries []models.Entry
	nextURL := fmt.Sprintf("https://blog.hatena.ne.jp/%s/%s/atom/entry", hatenaID, blogID)

	for nextURL != "" {
		fmt.Printf("取得中: %s\n", nextURL)

		entries, next, err := fetchBlogEntriesPage(hatenaID, apiKey, nextURL)
		if err != nil {
			return nil, err
		}

		// 増分取得モードの場合、既存日付以前の記事が見つかったら停止
		if !allMode && !sinceDate.IsZero() {
			var newEntries []models.Entry
			foundOldEntry := false

			for _, entry := range entries {
				if publishedTime, err := time.Parse(time.RFC3339, entry.Published); err == nil {
					// sinceDateより新しい記事のみを追加
					if publishedTime.After(sinceDate) {
						newEntries = append(newEntries, entry)
					} else {
						// 既存日付以前の記事が見つかったので停止
						foundOldEntry = true
						break
					}
				}
			}

			allEntries = append(allEntries, newEntries...)
			fmt.Printf("  %d件取得（新規: %d件、累計: %d件）\n", len(entries), len(newEntries), len(allEntries))

			// 古い記事に到達したら取得を停止
			if foundOldEntry {
				fmt.Printf("既存の記事に到達したため取得を停止します\n")
				break
			}
		} else {
			// 全記事取得モードまたは初回実行時
			allEntries = append(allEntries, entries...)
			fmt.Printf("  %d件取得（累計: %d件）\n", len(entries), len(allEntries))
		}

		nextURL = next
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

// getLatestEntryDate はentriesフォルダから最新の記事日付を取得する
func getLatestEntryDate() (time.Time, error) {
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
