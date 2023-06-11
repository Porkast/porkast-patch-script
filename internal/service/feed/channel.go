package feed

import (
	"context"
	"guoshao-fm-patch/internal/model/entity"

	"github.com/gogf/gf/v2/frame/g"
)

func PatchFeedChannelAuthor(ctx context.Context) {
	var (
		feedChannelList []entity.FeedChannel
		err             error
	)

	err = g.Model("feed_channel fc").Scan(&feedChannelList)
	if err != nil {
		g.Log().Line().Error(ctx, err)
		return
	}

	g.Log().Line().Infof(ctx, "rss feed channel total count : %d", len(feedChannelList))

	for _, channelInfo := range feedChannelList {
		var origAuthor = channelInfo.Author
		channelInfo.Author = formatFeedAuthor(channelInfo.Author)
		if origAuthor != channelInfo.Author {
			g.Log().Line().Infof(ctx, "Patch feed channel %s, from %s to %s", channelInfo.Title, origAuthor, channelInfo.Author)
			g.Model("feed_channel fc").Update(g.Map{"author": channelInfo.Author}, "id", channelInfo.Id)
		}
	}

}
