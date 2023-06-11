package cmd

import (
	"context"

	"github.com/gogf/gf/v2/os/gcmd"
)

var (
	Main = gcmd.Command{
		Name:  "main",
		Usage: "main",
		Brief: "start guoshao fm patch",
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {

			return nil
		},
	}
)
