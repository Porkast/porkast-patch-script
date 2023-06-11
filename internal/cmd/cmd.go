package cmd

import (
	"context"
	"guoshao-fm-patch/internal/service/feed"
	"os"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcmd"
	"github.com/gogf/gf/v2/os/genv"
)

var (
	Main = gcmd.Command{
		Name:  "main",
		Usage: "main",
		Brief: "start guoshao fm patch",
	}

	FeedChannelAuthorPatch = gcmd.Command{
		Name:  "FeedChannelAuthorPatch",
		Usage: "main",
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
		Usage: "main",
		Brief: "start guoshao fm FeedItemAuthorPatch",
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
            initConfig()
            g.Log().Line().Debug(ctx, "start guoshao fm FeedItemAuthorPatch")
            feed.PatchFeedItemAuthor(ctx)
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
