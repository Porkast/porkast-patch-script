// =================================================================================
// Code generated by GoFrame CLI tool. DO NOT EDIT. Created at 2023-08-12 16:28:28
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// KeywordSubscription is the golang structure for table keyword_subscription.
type KeywordSubscription struct {
	Id            string      `json:"id"            ` //
	Keyword       string      `json:"keyword"       ` //
	FeedChannelId string      `json:"feedChannelId" ` //
	FeedItemId    string      `json:"feedItemId"    ` //
	CreateTime    *gtime.Time `json:"createTime"    ` //
	OrderByDate   int         `json:"orderByDate"   ` //
	Lang          string      `json:"lang"          ` // feed language
}