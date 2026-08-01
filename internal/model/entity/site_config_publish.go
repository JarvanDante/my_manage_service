// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
package entity

import "github.com/gogf/gf/v2/os/gtime"

type SiteConfigPublish struct {
	Id        int64       `json:"id"        orm:"id"`
	SiteId    int64       `json:"siteId"    orm:"site_id"`
	Env       string      `json:"env"       orm:"env"`
	Version   int         `json:"version"   orm:"version"`
	Content   string      `json:"content"   orm:"content"`
	Operator  int64       `json:"operator"  orm:"operator"`
	Remark    string      `json:"remark"    orm:"remark"`
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at"`
}
