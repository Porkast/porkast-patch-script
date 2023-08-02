package search

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/gclient"
)

type Client struct {
	Ctx        context.Context
	httpClient *gclient.Client
}

var searchClient Client
var baseUrl string

func InitClient(ctx context.Context) {
	searchClient = Client{
		httpClient: gclient.New(),
	}

	username := os.Getenv("ZINC_FIRST_ADMIN_USER")
	password := os.Getenv("ZINC_FIRST_ADMIN_PASSWORD")
	baseUrl = os.Getenv("ZINC_BASE_URL")
	if username == "" || password == "" {
		panic("ZINC_FIRST_ADMIN_USER and ZINC_FIRST_ADMIN_PASSWORD must be set")
	}
	g.Log().Line().Debugf(ctx, "zinc init with user: %s pass: %s", username, password)
	searchClient.httpClient.SetBasicAuth(username, password)
	searchClient.httpClient.SetHeader("Content-Type", "application/json")
}

func GetClient(ctx context.Context) Client {
	searchClient.Ctx = ctx
	return searchClient
}

func (c Client) doPost(ctx context.Context, url, body string) (err error) {

	resp, err := c.httpClient.Post(ctx, url, body)
	defer func(resp *gclient.Response) {
		if rec := recover(); rec != nil {
			if resp == nil {
				g.Log().Line().Error(ctx, fmt.Sprintf("do post failed: \n%s\n", rec))
			} else {
				g.Log().Line().Error(ctx, fmt.Sprintf("do post failed with response %s: \n%s\n", resp.ReadAllString(), rec))
			}
		}
		if resp != nil {
			resp.Close()
		}
	}(resp)

	if err != nil {
		g.Log().Line().Error(ctx, err)
	} else if resp != nil && resp.StatusCode != 200 {
		err = errors.New(resp.ReadAllString())
		g.Log().Line().Error(ctx, err)
	}

	return
}
