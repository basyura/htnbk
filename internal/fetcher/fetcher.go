package fetcher

import (
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"time"

	"htnblg-export/internal/models"
)

// FetchAllBlogEntries はブログの記事を効率的に取得する
// allMode が true の場合はすべての記事を取得
// allMode が false の場合は sinceDate 以降の記事のみを取得
func FetchAllBlogEntries(hatenaID, blogID, apiKey string, allMode bool, sinceDate time.Time) ([]models.Entry, error) {
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

// fetchBlogEntriesPage は単一ページの記事を取得する
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
