package models

type Link struct {
	Rel  string `xml:"rel,attr"`
	Href string `xml:"href,attr"`
}