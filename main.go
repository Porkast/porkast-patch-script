package main

import (
	_ "guoshao-fm-patch/internal/packed"

	"github.com/gogf/gf/v2/os/gctx"
	"guoshao-fm-patch/internal/cmd"
)

func main() {
	cmd.Main.Run(gctx.New())
}
