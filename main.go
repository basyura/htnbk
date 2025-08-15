package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"htnblg-export/internal/fetcher"
	"htnblg-export/internal/storage"
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
		sinceDate, err = storage.GetLatestEntryDate()
		if err != nil {
			fmt.Printf("警告: 最新記事日付の取得に失敗: %v\n", err)
		} else if !sinceDate.IsZero() {
			fmt.Printf("最新記事日付: %s\n", sinceDate.Format("2006-01-02"))
		} else {
			fmt.Printf("初回実行のため、すべての記事を取得します。\n")
		}
	} else {
		fmt.Printf("--all モードですべての記事を取得します。\n")
	}

	// 効率的な取得：APIレベルで必要な分だけ取得
	entries, err := fetcher.FetchAllBlogEntries(hatenaID, blogID, apiKey, allMode, sinceDate)
	if err != nil {
		return err
	}

	// 取得結果の表示
	if !allMode && !sinceDate.IsZero() {
		fmt.Printf("新しい記事: %d件\n", len(entries))
	} else {
		fmt.Printf("総記事数: %d件\n", len(entries))
	}

	if len(entries) == 0 {
		fmt.Println("新しい記事はありません。")
		return nil
	}

	fmt.Println(strings.Repeat("-", 50))

	for i, entry := range entries {
		publishedTime, err := time.Parse(time.RFC3339, entry.Published)
		if err != nil {
			fmt.Printf("日付解析エラー: %v\n", err)
			continue
		}
		fmt.Printf("%04d : %s - %s\n", i+1, publishedTime.Format("2006-01-02"), entry.Title)
		fmt.Printf("     ID: %s\n", entry.ID)

		// ファイルパスを生成（年/月/ファイル名）
		filePath, err := storage.GenerateFilePath(entry.Published, entry.Title)
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
		err = storage.SaveEntryToFile(filePath, &entry)
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
