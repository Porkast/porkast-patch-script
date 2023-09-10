package feed

import (
	"context"
	"guoshao-fm-patch/internal/consts"
	"guoshao-fm-patch/internal/model/entity"
	"guoshao-fm-patch/internal/service/cache"
	"guoshao-fm-patch/internal/service/internal/dao"
	"guoshao-fm-patch/internal/service/search"
	"strconv"
	"sync"

	"github.com/anaskhan96/soup"
	"github.com/gogf/gf/v2/encoding/ghash"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/grpool"
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
	err = cache.SetCache(ctx, gconv.String(consts.FEED_CHANNEL_TOTAL_COUNT), gconv.String(totalCount), int(24*60*60))
	if err != nil {
		panic(err)
	}
	return
}

func SetFeedChannelToZincsearch(ctx context.Context, feedChannel entity.FeedChannel) {
	esFeedChannel := entity.FeedChannelESData{}
	gconv.Struct(feedChannel, &esFeedChannel)
	rootDocs := soup.HTMLParse(esFeedChannel.ChannelDesc)
	esFeedChannel.TextChannelDesc = rootDocs.FullText()
	search.GetClient(ctx).InsertFeedChannel(esFeedChannel)
}

func FilterDuplicatedChannelInfo(ctx context.Context) {
	var (
		err                   error
		totalCount            int
		offset                int
		limit                 = 1000
		channelInfoList       []entity.FeedChannel
		channelInfoMap        map[string]entity.FeedChannel
		duplicatedChannelList []DuplicatedChannelInfo
	)

	totalCount, err = dao.GetZHFeedChannelTotalCount(ctx)
	if err != nil {
		panic(err)
	}
	g.Log().Line().Infof(ctx, "The channel info total count is %d", totalCount)
	channelInfoMap = make(map[string]entity.FeedChannel)
	duplicatedChannelList = make([]DuplicatedChannelInfo, 0)
	for offset < totalCount {
		channelInfoList, err = dao.GetChannelList(ctx, offset, limit)
		offset = offset + limit
		for _, channleInfoItem := range channelInfoList {
			key := channleInfoItem.Link + channleInfoItem.Title
			existChannelInfo := channelInfoMap[key]
			if existChannelInfo.Id == "" {
				channelInfoMap[key] = channleInfoItem
			} else {
				duplicatedInfo := DuplicatedChannelInfo{
					Id:        channleInfoItem.Id,
					Id2:       existChannelInfo.Id,
					Link:      channleInfoItem.Link,
					Link2:     existChannelInfo.Link,
					Title:     channleInfoItem.Title,
					Title2:    existChannelInfo.Title,
					FeedLink:  channleInfoItem.FeedLink,
					FeedLink2: existChannelInfo.FeedLink,
				}
				duplicatedChannelList = append(duplicatedChannelList, duplicatedInfo)
			}
		}
	}

	resultJsonStr := gjson.MustEncodeString(duplicatedChannelList)
	gfile.PutContents("./duplicated_channle_info.json", resultJsonStr)
	g.Log().Line().Infof(ctx, "The duplicated channel count is %d", len(duplicatedChannelList))
}

func MigrateFeedChannelAndItem(ctx context.Context) {
	var (
		err error
		// mu              = &sync.Mutex{}
		totalCount      int
		offset          int
		limit           = 1000
		channelInfoList []entity.FeedChannel
	)

	totalCount, err = dao.GetFeedChannelTotalCount(ctx)
	if err != nil {
		panic(err)
	}
	g.Log().Line().Infof(ctx, "The channel info total count is %d", totalCount)
	for offset < totalCount {
		// feedChannelMigrateList := make([]entity.FeedChannelMigrate, 0)
		channelInfoList, err = dao.GetChannelList(ctx, offset, limit)
		g.Log().Line().Infof(ctx, "start from offset %d", offset)
		offset = offset + limit
		wg := sync.WaitGroup{}
		pool := grpool.New(10)
		for _, channelInfoItem := range channelInfoList {
			wg.Add(1)
			feedChannelTemp := channelInfoItem
			pool.Add(ctx, func(ctx context.Context) {
				defer wg.Done()
				if feedChannelTemp.FeedLink == "" {
					g.Log().Line().Infof(ctx, "The channel %s feed link is empty", feedChannelTemp.Id)
					return
				}
				newId := strconv.FormatUint(ghash.RS64([]byte(feedChannelTemp.FeedLink+feedChannelTemp.Title)), 32)
				feedChannelMigrate := entity.FeedChannelMigrate{}
				gconv.Struct(feedChannelTemp, &feedChannelMigrate)
				feedChannelMigrate.Ido = feedChannelTemp.Id
				feedChannelMigrate.Id = newId
				migrateFeedItem(ctx, feedChannelTemp.Id, newId)
				// mu.Lock()
				// feedChannelMigrateList = append(feedChannelMigrateList, feedChannelMigrate)
				// mu.Unlock()
			})
		}

		//do insert or update MigrateFeedChannel
		wg.Wait()
		// dao.FeedChannelMigrate.Ctx(ctx).Data(feedChannelMigrateList).Save()
	}

}

func migrateFeedItem(ctx context.Context, originalChannelId, newChannelId string) (err error) {

	feedItemMigrateList := make([]entity.FeedItemMigrate, 0)
	feedItemList, err := dao.GetFeedItemsByChannelId(ctx, originalChannelId)
	if len(feedItemList) == 0 {
		g.Log().Line().Infof(ctx, "Feed channel %s item is empty, ignore update feed_item", originalChannelId)
		return
	}

	for _, feedItem := range feedItemList {
		feedItemMigrate := entity.FeedItemMigrate{}
		gconv.Struct(feedItem, &feedItemMigrate)
		feedItemMigrate.ChannelId = newChannelId
		var itemID = strconv.FormatUint(ghash.RS64([]byte(newChannelId+feedItemMigrate.Title)), 32)
		feedItemMigrate.Id = itemID
		feedItemMigrateList = append(feedItemMigrateList, feedItemMigrate)
	}

	dao.FeedItemMigrate.Ctx(ctx).Data(feedItemMigrateList).Save()

	return
}
