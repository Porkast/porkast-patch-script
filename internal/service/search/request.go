package search

import (
	"errors"
	"fmt"
	"guoshao-fm-patch/internal/model/entity"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/gclient"
)

func (c Client) CreateIndex( body string) (err error) {

	apiUrl := fmt.Sprintf(baseUrl + "/api/index")
	resp, err := c.httpClient.Post(c.Ctx, apiUrl, body)
	defer func(resp *gclient.Response) {
		if rec := recover(); rec != nil {
			if resp == nil {
				g.Log().Line().Error(c.Ctx, fmt.Sprintf("Create index failed: \n%s\n", rec))
			} else {
				g.Log().Line().Error(c.Ctx, fmt.Sprintf("Create index failed with response %s: \n%s\n", resp.ReadAllString(), rec))
			}
		}
		if resp != nil {
			resp.Close()
		}
	}(resp)

	if err != nil {
		g.Log().Line().Error(c.Ctx, err)
	} else if resp != nil && resp.StatusCode != 200 {
		err = errors.New(resp.ReadAllString())
		g.Log().Line().Error(c.Ctx, err)
	}

	return
}

func (c Client) InsertFeedChannel(searchChannelModel entity.FeedChannelESData) (err error) {

	apiUrl := baseUrl + "/api/feed_channel/_doc"
	body := gjson.MustEncodeString(searchChannelModel)
	err = c.doPost(c.Ctx, apiUrl, body)
	if err != nil {
		g.Log().Line().Error(c.Ctx, fmt.Sprintf("Insert feed channel document failed: \n%s\n", err))
	}

	return
}

func (c Client) InsertFeedItemBulk(itemBulks FeedItemBulk) (err error) {

	apiUrl := baseUrl + "/api/_bulkv2"
	body := gjson.MustEncodeString(itemBulks)
	err = c.doPost(c.Ctx, apiUrl, body)
	if err != nil {
		g.Log().Line().Error(c.Ctx, fmt.Sprintf("Bulk insert feed item document failed: \n%s\n", err))
	}

	return
}
