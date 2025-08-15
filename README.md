# htnblg-export

はてなブログエクスポートツール - Hatena Blog Export Tool

## 概要

はてなブログから記事を抜き出して、エントリーごとにファイルに出力するバックアップ用ツールです。はてなブログのAtomPub APIを使用して、ブログの全記事を年月ごとのディレクトリに整理して保存します。

## 特徴

- **完全バックアップ**: ブログの全記事を取得
- **自動整理**: 年/月のディレクトリ構造で自動整理
- **ページネーション対応**: 大量の記事も段階的に取得
- **詳細情報表示**: 記事数、公開日、タイトルを表示
- **ファイル出力**: 各記事を個別のファイルとして保存

## 使用方法

### 基本的な使用法

```bash
./htnblg-export <はてなID> <ブログID> <APIキー>
```

### 使用例

```bash
./htnblg-export basyura blog.basyura.org your_api_key_here
```

### パラメータ

- **はてなID**: あなたのはてなID
- **ブログID**: ブログのドメイン名（例: blog.example.com）
- **APIキー**: はてなブログのAPIキー

## セットアップ

### 1. はてなブログAPIキーの取得

1. [はてなブログの設定画面](https://blog.hatena.ne.jp/my/config)にアクセス
2. 「詳細設定」タブを選択
3. 「AtomPub」セクションでAPIキーを確認・生成

### 2. ブログIDの確認

ブログIDは、あなたのブログのドメイン名です：
- 独自ドメインの場合: `blog.example.com`
- はてなサブドメインの場合: `username.hatenablog.com`

## 出力形式

### ディレクトリ構造

```
output/
├── 2024/
│   ├── 01/  # 1月の記事
│   ├── 02/  # 2月の記事
│   └── ...
└── 2023/
    ├── 01/
    └── ...
```

### ファイル命名規則

各記事は以下の形式でファイル名が付けられます：
```
YYYY-MM-DD_記事タイトル.xml
```

## ライセンス

このプロジェクトは MIT ライセンスの下で公開されています。詳細については [LICENSE](LICENSE) ファイルを参照してください。

## 関連リンク

- [はてなブログ AtomPub API ドキュメント](https://developer.hatena.ne.jp/ja/documents/blog/apis/atom/)
- [はてなブログ設定画面](https://blog.hatena.ne.jp/my/config)
