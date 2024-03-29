package feed

import (
	"context"
	"guoshao-fm-patch/internal/consts"
	"guoshao-fm-patch/internal/dto"
	"guoshao-fm-patch/internal/model/entity"
	"guoshao-fm-patch/internal/service/cache"
	"guoshao-fm-patch/internal/service/internal/dao"
	"guoshao-fm-patch/internal/service/search"
	"strconv"
	"sync"
	"sync/atomic"

	"github.com/anaskhan96/soup"
	"github.com/gogf/gf/v2/encoding/ghash"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/grpool"
	"github.com/gogf/gf/v2/os/gtime"
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
	err = cache.SetCache(ctx, gconv.String(consts.FEED_ITEM_TOTAL_COUNT), gconv.String(totalCount), int(24*60*60))
	if err != nil {
		panic(err)
	}
	return
}

func SetLatestFeedItems(ctx context.Context) (err error) {
	var (
		startDate    *gtime.Time
		startDateStr string
		endDate      *gtime.Time
		endDateStr   string
		itemList     []dto.FeedItem
		itemListJson *gjson.Json
	)

	startDate = gtime.Now().StartOfDay()
	endDate = gtime.Now().EndOfDay()

	startDateStr = startDate.ISO8601()
	endDateStr = endDate.ISO8601()

	itemList = dao.GetFeedItemListByPubDate(ctx, startDateStr, endDateStr)
	if err != nil {
		g.Log().Line().Error(ctx, "Get latest feed items failed: ", err)
		return
	}

	if len(itemList) == 0 {
		return
	}

	itemListJson = gjson.New(itemList)
	if err != nil {
		g.Log().Line().Error(ctx, "Decode feed item list to json failed", err)
		return
	}
	cache.SetCache(ctx, gconv.String(consts.TODAY_FEED_ITEM_LIST), itemListJson.MustToJsonString(), 0)

	return
}

func SetFeedItemsToZincsearch(ctx context.Context, feedChannel entity.FeedChannel, feedItemList []entity.FeedItem) {

	if len(feedItemList) == 0 {
		return
	}
	esFeedItemList := make([]entity.FeedItemESData, 0)
	for _, feedItem := range feedItemList {
		esFeedItem := entity.FeedItemESData{}
		gconv.Struct(feedItem, &esFeedItem)
		if esFeedItem.Author == "" {
			esFeedItem.Author = feedChannel.Author
		}
		rootDocs := soup.HTMLParse(esFeedItem.Description)
		esFeedItem.TextDescription = rootDocs.FullText()
		esFeedItem.ChannelImageUrl = feedChannel.ImageUrl
		esFeedItem.ChannelTitle = feedChannel.Title
		esFeedItem.FeedLink = feedChannel.FeedLink
		esFeedItem.Language = feedChannel.Language
		esFeedItem.Ids = feedItem.Id
		esFeedItemList = append(esFeedItemList, esFeedItem)
	}

	bulk := search.FeedItemBulk{
		Index:   "feed_item",
		Records: esFeedItemList,
	}

	search.GetClient(ctx).InsertFeedItemBulk(bulk)
}

func GetMigrateFeedItemTotalCount(ctx context.Context) {
	var (
		itemTotalCount int64
	)

	channelInfoList, _ := dao.GetMigrateZHFeedChannelList(ctx)
	wg := sync.WaitGroup{}
	pool := grpool.New(10)
	for _, channelItem := range channelInfoList {
		wg.Add(1)
		feedChannelTemp := channelItem
		pool.Add(ctx, func(ctx context.Context) {
			defer wg.Done()
			count, err := dao.GetMigrateFeedItemCountByChannelId(ctx, feedChannelTemp.Id)
			if err != nil {
				g.Log().Line().Errorf(ctx, "Get channel %s item total count failed :\n%d", feedChannelTemp.Title, err)
				return
			}

			g.Log().Line().Infof(ctx, "The channel %s items total count is %d", feedChannelTemp.Id, count)
			atomic.AddInt64(&itemTotalCount, gconv.Int64(count))
		})
	}
	wg.Wait()
	g.Log().Line().Infof(ctx, "The all migrate items total count is %d", itemTotalCount)
}

func MigrateListenLaterFeedItem(ctx context.Context) {

	userListenLaterEntities, err := dao.GetAllListenLaterEntities(ctx)
	if err != nil {
		panic(err)
	}

	for _, userListenLater := range userListenLaterEntities {
		channelId := userListenLater.ChannelId
		itemId := userListenLater.ItemId
		itemInfo, err := dao.GetFeedItemByChannelIdAndItemId(ctx, channelId, itemId)
		if err != nil {
			g.Log().Line().Errorf(ctx, "get channel item info by channel id %s and item id %s faild:\n%s", channelId, itemId, err)
			continue
		}

		var newItemId = strconv.FormatUint(ghash.RS64([]byte(channelId+itemInfo.Title)), 32)
		userListenLater.ItemId = newItemId
		dao.UserListenLater.Ctx(ctx).Save(userListenLater)
	}

}
