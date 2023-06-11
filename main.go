package main

import (
	"guoshao-fm-patch/internal/cmd"
	_ "guoshao-fm-patch/internal/packed"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
)

func main() {
	err := cmd.Main.AddCommand(&cmd.FeedChannelAuthorPatch, &cmd.FeedItemAuthorPatch)
	if err != nil {
		g.Log().Line().Fatal(gctx.New(), err)
	}

	cmd.Main.Run(gctx.New())
}
