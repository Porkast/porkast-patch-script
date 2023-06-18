package feed

import (
	"context"
	"guoshao-fm-patch/internal/consts"
	"guoshao-fm-patch/internal/model/entity"
	"guoshao-fm-patch/internal/service/cache"
	"guoshao-fm-patch/internal/service/internal/dao"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
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

func SetZHChannelTotalCountToCache(ctx context.Context) (err error) {
	var (
		totalCount int
	)

	totalCount, err = dao.GetZHFeedChannelTotalCount(ctx)
	if err != nil {
		g.Log().Line().Error(ctx, "Get feed channel total count failed : ", err)
		return
	}

	g.Log().Line().Info(ctx, "The all ZH channel total count is ", totalCount)
	cache.SetCache(ctx, gconv.String(consts.FEED_CHANNEL_TOTAL_COUNT), gconv.String(totalCount), 0)
	return
}
