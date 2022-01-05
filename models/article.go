package models

import "html/template"

type Article struct {
	ID           string        `json:"id" bson:"_id,omitempty"`
	Title        string        `json:"title"`
	Subtitle     string        `json:"subtitle"`
	ContentMD    string        `json:"contentMD"`
	Content      template.HTML `json:"content"`
	CreatedDate  int64         `json:"createdDate"`
	ModifiedDate int64         `json:"modifiedDate"`
	Tags         []string      `json:"tags"`
	Categories   []string      `json:"categories"`
	Author       []string      `json:"author"`
	NeedsTOC     bool          `json:"needsTOC"`
	TOC          string        `json:"toc"`
	IsInSeries   bool          `json:"isInSeries"`
	Series       string        `json:"series"`
}
