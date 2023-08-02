package entity

import "github.com/gogf/gf/v2/os/gtime"

type FeedItemESData struct {
	Id              string      `json:"id"              ` //
	ChannelId       string      `json:"channelId"       ` //
	Title           string      `json:"title"           ` //
	Link            string      `json:"link"            ` //
	PubDate         *gtime.Time `json:"pubDate"         ` //
	Author          string      `json:"author"          ` //
	InputDate       *gtime.Time `json:"inputDate"       ` //
	ImageUrl        string      `json:"imageUrl"        ` //
	EnclosureUrl    string      `json:"enclosureUrl"    ` //
	EnclosureType   string      `json:"enclosureType"   ` //
	EnclosureLength string      `json:"enclosureLength" ` //
	Duration        string      `json:"duration"        ` //
	Episode         string      `json:"episode"         ` //
	Explicit        string      `json:"explicit"        ` //
	Season          string      `json:"season"          ` //
	EpisodeType     string      `json:"episodeType"     ` //
	Description     string      `json:"description"     ` //
	TextDescription string      `json:"textDescription" ` //
	ChannelImageUrl string      `json:"channelImageUrl" ` //
	ChannelTitle    string      `json:"channelTitle"    ` //
	FeedLink        string      `json:"feedLink"        ` //
	Language        string      `json:"language"        ` //
}
