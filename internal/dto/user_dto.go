package dto

import "github.com/gogf/gf/v2/os/gtime"

type UserInfo struct {
	Id         string      `json:"id"         ` //
	Username   string      `json:"username"   ` //
	Nickname   string      `json:"nickname"   ` //
	Token      string      `json:"token"`
	Password   string      `json:"password"   ` //
	Email      string      `json:"email"      ` //
	Phone      string      `json:"phone"      ` //
	RegDate    *gtime.Time `json:"regDate"    ` //
	UpdateDate *gtime.Time `json:"updateDate" ` //
}
