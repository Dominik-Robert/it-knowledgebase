package models

type Article struct {
	Title        string   `json:"title"`
	Subtitle     string   `json:"subtitle"`
	ContentMD    string   `json:"contentMD"`
	Content      string   `json:"content"`
	CreatedDate  string   `json:"createdDate"`
	ModifiedDate string   `json:"modifiedDate"`
	Tags         []string `json:"tags"`
	Categories   []string `json:"categories"`
	Author       []string `json:"author"`
	NeedsTOC     bool     `json:"needsTOC"`
	TOC          string   `json:"toc"`
	IsInSeries   bool     `json:"isInSeries"`
	Series       string   `json:"series"`
}
