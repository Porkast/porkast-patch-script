package cmd

import (
	"context"
	"guoshao-fm-patch/internal/model/entity"
	"guoshao-fm-patch/internal/service/cache"
	"guoshao-fm-patch/internal/service/feed"
	"guoshao-fm-patch/internal/service/search"
	"os"
	"sync"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcmd"
	"github.com/gogf/gf/v2/os/genv"
	"github.com/gogf/gf/v2/os/grpool"
)

var (
	Main = gcmd.Command{
		Name:  "main",
		Usage: "main",
		Brief: "start guoshao fm patch",
	}

	FeedChannelAuthorPatch = gcmd.Command{
		Name:  "FeedChannelAuthorPatch",
		Usage: "patch",
		Brief: "start guoshao fm FeedChannelAuthorPatch",
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			initConfig()
			g.Log().Line().Debug(ctx, "start guoshao fm FeedChannelAuthorPatch")
			feed.PatchFeedChannelAuthor(ctx)
			return nil
		},
	}

	FeedItemAuthorPatch = gcmd.Command{
		Name:  "FeedItemAuthorPatch",
		Usage: "patch",
		Brief: "start guoshao fm FeedItemAuthorPatch",
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			initConfig()
			g.Log().Line().Debug(ctx, "start guoshao fm FeedItemAuthorPatch")
			feed.PatchFeedItemAuthor(ctx)
			return nil
		},
	}

	FeedChannelItemTotalCountPatch = gcmd.Command{
		Name:  "FeedChannelItemTotalCountPatch",
		Usage: "patch",
		Brief: "start guoshao fm feed total count patch",
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			initConfig()
			cache.InitCache(ctx)
			g.Log().Line().Debug(ctx, "start guoshao fm feed total count patch")
			feed.SetZHChannelTotalCountToCache(ctx)
			feed.SetZHItemTotalCountToCache(ctx)
			return nil
		},
	}

	SetLatestItemToCachePatch = gcmd.Command{
		Name:  "SetLatestItemToCachePatch",
		Usage: "patch",
		Brief: "start guoshao fm SetLatestItemToCachePatch patch",
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			initConfig()
			cache.InitCache(ctx)
			feed.SetLatestFeedItems(ctx)
			g.Log().Line().Debug(ctx, "start guoshao fm SetLatestItemToCachePatch patch")
			return nil
		},
	}

	AddZincsearchIndex = gcmd.Command{
		Name:  "AddZincsearchIndex",
		Usage: "patch",
		Brief: "start guoshao fm AddZincsearchIndex patch",
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			initConfig()
			search.InitClient(ctx)
			search.GetClient(ctx).CreateIndex(search.CREATE_FEED_ITEM_INDEX_REQUEST)
			search.GetClient(ctx).CreateIndex(search.CREATE_FEED_CHANNEL_INDEX_REQUEST)
			g.Log().Line().Debug(ctx, "start guoshao fm AddZincsearchIndex patch")
			return nil
		},
	}

	PatchZincsearchData = gcmd.Command{
		Name:  "PatchZincsearchData",
		Usage: "patch",
		Brief: "start guoshao fm PatchZincsearchData",
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			initConfig()
			search.InitClient(ctx)
			g.Log().Line().Debug(ctx, "start guoshao fm PatchZincsearchData patch")

			feedChannelList := make([]entity.FeedChannel, 0)
			err = g.Model("feed_channel fc").Scan(&feedChannelList)
			if err != nil {
				g.Log().Line().Error(ctx, err)
				return
			}

			g.Log().Line().Infof(ctx, "rss feed channel total count : %d", len(feedChannelList))
			wg := sync.WaitGroup{}
			pool := grpool.New(10)
			for _, feedChannel := range feedChannelList {
				wg.Add(1)
				feedChannelTemp := feedChannel
				pool.Add(ctx, func(ctx context.Context) {
					defer wg.Done()
					feedItemList := make([]entity.FeedItem, 0)
					err := g.Model("feed_item fi").
						Where("fi.channel_id=?", feedChannelTemp.Id).
						Scan(&feedItemList)
					if err != nil {
						g.Log().Line().Error(ctx, err)
					}
					feed.SetFeedChannelToZincsearch(ctx, feedChannelTemp)
					feed.SetFeedItemsToZincsearch(ctx, feedChannelTemp, feedItemList)
					g.Log().Line().Infof(ctx, "channel %s total item count : %d", feedChannelTemp.Title, len(feedItemList))
				})
			}
			g.Log().Line().Infof(ctx, "start patch with pool workers %d, jobs %d ", pool.Size(), pool.Jobs())
			wg.Wait()
			return nil
		},
	}

	DuplicatedChannelPatch = gcmd.Command{
		Name:  "DuplicatedChannelPatch",
		Usage: "patch",
		Brief: "start guoshao fm DuplicatedChannelPatch patch",
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			initConfig()
			g.Log().Line().Debug(ctx, "start guoshao fm DuplicatedChannelPatch patch")
			g.Log().Line().Debug(ctx, "db host %s", g.DB().GetConfig().Host)
			feed.FilterDuplicatedChannelInfo(ctx)
			return nil
		},
	}

	MigrateFeedChannelAndItemDBTable = gcmd.Command{
		Name:  "MigrateFeedChannelAndItemDBTable",
		Usage: "patch",
		Brief: "start guoshao fm MigrateFeedChannelAndItemDBTable patch",
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			initConfig()
			g.Log().Line().Debug(ctx, "start guoshao fm MigrateFeedChannelAndItemDBTable patch")
			g.Log().Line().Debug(ctx, "db host %s", g.DB().GetConfig().Host)
			feed.MigrateFeedChannelAndItem(ctx)
			return nil
		},
	}
)

func initConfig() {
	if os.Getenv("env") == "dev" {
		genv.Set("GF_GCFG_FILE", "config.dev.yaml")
	} else if os.Getenv("env") == "prod" {
		genv.Set("GF_GCFG_FILE", "config.prod.yaml")
	} else {
		genv.Set("GF_GCFG_FILE", "config.yaml")
	}
	g.I18n().SetPath("./resource/i18n")
}
