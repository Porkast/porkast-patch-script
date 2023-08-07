package feed

import (
	"context"
	"guoshao-fm-patch/internal/consts"
	"guoshao-fm-patch/internal/model/entity"
	"guoshao-fm-patch/internal/service/cache"
	"guoshao-fm-patch/internal/service/internal/dao"
	"guoshao-fm-patch/internal/service/search"

	"github.com/anaskhan96/soup"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gfile"
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
