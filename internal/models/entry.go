package models

type Entry struct {
	ID         string     `xml:"id"`
	Title      string     `xml:"title"`
	Published  string     `xml:"published"`
	Updated    string     `xml:"updated"`
	Edited     string     `xml:"http://www.w3.org/2007/app edited"`
	Author     Author     `xml:"author"`
	Links      []Link     `xml:"link"`
	Categories []Category `xml:"category"`
	Content    string     `xml:"content"`
	Control    AppControl `xml:"http://www.w3.org/2007/app control"`
	CustomURL  string     `xml:"http://www.hatena.ne.jp/info/xmlns#hatenablog custom-url"`
}
