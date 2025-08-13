package main

import (
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"os"
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

	err := fetchBlogEntries(hatenaID, blogID, apiKey)
	if err != nil {
		fmt.Printf("エラー: %v\n", err)
		os.Exit(1)
	}
}

func fetchBlogEntries(hatenaID, blogID, apiKey string) error {
	url := fmt.Sprintf("https://blog.hatena.ne.jp/%s/%s/atom/entry", hatenaID, blogID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	// Basic認証のヘッダーを設定
	auth := base64.StdEncoding.EncodeToString([]byte(hatenaID + ":" + apiKey))
	req.Header.Set("Authorization", "Basic "+auth)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTPエラー: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var feed Feed
	err = xml.Unmarshal(body, &feed)
	if err != nil {
		return err
	}

	for i, entry := range feed.Entries {
		publishedTime, err := time.Parse(time.RFC3339, entry.Published)
		if err != nil {
			fmt.Printf("日付解析エラー: %v\n", err)
			continue
		}
		
		fmt.Printf("%04d : %s - %s\n", i+1, publishedTime.Format("2006-01-02"), entry.Title)
	}

	return nil
}
