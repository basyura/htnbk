package models

type Link struct {
	Rel  string `xml:"rel,attr"`
	Href string `xml:"href,attr"`
}

type Category struct {
	Term string `xml:"term,attr"`
}

type Author struct {
	Name string `xml:"name"`
}

type AppControl struct {
	Draft   string `xml:"http://www.w3.org/2007/app draft"`
	Preview string `xml:"http://www.w3.org/2007/app preview"`
}
