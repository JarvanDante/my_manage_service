// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
package entity

import "github.com/gogf/gf/v2/os/gtime"

type SiteDomain struct {
	Id        int64       `json:"id"        orm:"id"`
	SiteId    int64       `json:"siteId"    orm:"site_id"`
	Domain    string      `json:"domain"    orm:"domain"`
	IsMain    int         `json:"isMain"    orm:"is_main"`
	Https     int         `json:"https"     orm:"https"`
	Status    int         `json:"status"    orm:"status"`
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at"`
}
