package feed

import (
	"context"
	"guoshao-fm-patch/internal/consts"
	"guoshao-fm-patch/internal/model/entity"
	"guoshao-fm-patch/internal/service/cache"
	"guoshao-fm-patch/internal/service/internal/dao"
	"sync"
	"sync/atomic"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/grpool"
	"github.com/gogf/gf/v2/util/gconv"
)

func PatchFeedItemAuthor(ctx context.Context) {
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

	wg := sync.WaitGroup{}
	pool := grpool.New(100)
	for _, channelInfo := range feedChannelList {
		var (
			channelId string
			itemList  []entity.FeedItem
		)
		wg.Add(1)

		channelInfoTemp := channelInfo
		pool.Add(ctx, func(ctx context.Context) {
			channelId = channelInfoTemp.Id
			itemList, err = dao.GetFeedItemsByChannelId(ctx, channelId)
			for _, item := range itemList {
				if item.Author == "" {
					item.Author = channelInfoTemp.Author
				}
				if item.Author == "" && channelInfoTemp.OwnerName != "" {
					item.Author = channelInfoTemp.OwnerName
				}
				var origAuthor = item.Author
				item.Author = formatFeedAuthor(item.Author)
				if origAuthor != item.Author {
					g.Log().Line().Infof(ctx, "Patch feed channel %s item %s, from %s to %s", channelId, item.Id, origAuthor, item.Author)
					g.Model("feed_item").Update(g.Map{"author": item.Author}, "id", item.Id)
				}
			}
			wg.Done()

		})
	}
	g.Log().Line().Infof(ctx, "start patch with pool workers %d, jobs %d ", pool.Size(), pool.Jobs())
	wg.Wait()
}

func SetZHItemTotalCountToCache(ctx context.Context) (err error) {
	var (
		totalCount      int64
		feedChannelList []entity.FeedChannel
	)

	wg := sync.WaitGroup{}
	pool := grpool.New(100)
	feedChannelList, err = dao.GetZHFeedChannelList(ctx)
	for _, feedChannel := range feedChannelList {
		feedChannelTemp := feedChannel
		wg.Add(1)
		pool.Add(ctx, func(ctx context.Context) {
			defer wg.Done()
			count, err := dao.GetFeedItemCountByChannelId(ctx, feedChannelTemp.Id)
			if err == nil {
				g.Log().Line().Infof(ctx, "The channel %s item total count is %d", feedChannelTemp.Title, count)
				atomic.AddInt64(&totalCount, gconv.Int64(count))
			}
		})
	}

	wg.Wait()
	g.Log().Line().Infof(ctx, "The all ZH items total count is %d", totalCount)
	cache.SetCache(ctx, gconv.String(consts.FEED_ITEM_TOTAL_COUNT), gconv.String(totalCount), 0)
	return
}