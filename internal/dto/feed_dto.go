package dto

import "github.com/gogf/gf/v2/os/gtime"

type FeedItem struct {
	Id                    string
	ChannelId             string
	Title                 string
	HighlightTitle        string
	Link                  string
	PubDate               string
	Author                string
	InputDate             *gtime.Time
	ImageUrl              string
	EnclosureUrl          string
	EnclosureType         string
	EnclosureLength       string
	Duration              string
	Episode               string
	Explicit              string
	Season                string
	EpisodeType           string
	Description           string
	TextDescription       string
	ChannelImageUrl       string
	ChannelTitle          string
	HighlightChannelTitle string
	FeedLink              string
	Count                 int
	TookTime              float64
	HasThumbnail          bool
}

type FeedChannel struct {
	Id              string
	Title           string
	ChannelDesc     string
	TextChannelDesc string
	ImageUrl        string
	Link            string
	FeedLink        string
	Launguage       string
	FeedType        string
	Categories      []string
	Author          string
	OwnerName       string
	OwnerEmail      string
	Items           []FeedItem
	Count           int
	Copyright       string
	Language        string
	TookTime        float64
	HasThumbnail    bool
}
