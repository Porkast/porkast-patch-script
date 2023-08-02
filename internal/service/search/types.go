package search

import "guoshao-fm-patch/internal/model/entity"

type FeedItemBulk struct {
	Index   string                  `json:"index"`
	Records []entity.FeedItemESData `json:"records"`
}
