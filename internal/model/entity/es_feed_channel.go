package entity

// FeedItem is the golang structure for table feed_item.
type FeedChannelESData struct {
	Id              string `json:"id"          `     //
	Title           string `json:"title"       `     //
	ChannelDesc     string `json:"channelDesc" `     //
	TextChannelDesc string `json:"textChannelDesc" ` //
	ImageUrl        string `json:"imageUrl"    `     //
	Link            string `json:"link"        `     //
	FeedLink        string `json:"feedLink"    `     //
	Copyright       string `json:"copyright"   `     //
	Language        string `json:"language"    `     //
	Author          string `json:"author"      `     //
	OwnerName       string `json:"ownerName"   `     //
	OwnerEmail      string `json:"ownerEmail"  `     //
	FeedType        string `json:"feedType"    `     //
	Categories      string `json:"categories"  `     //
}
