// Code maintained manually (平台服务配置发布记录).
package entity

import "github.com/gogf/gf/v2/os/gtime"

type ServiceConfigPublish struct {
	Id        int64       `json:"id"        orm:"id"`
	Service   string      `json:"service"   orm:"service"`
	Env       string      `json:"env"       orm:"env"`
	Version   int         `json:"version"   orm:"version"`
	Content   string      `json:"content"   orm:"content"`
	Operator  int64       `json:"operator"  orm:"operator"`
	Remark    string      `json:"remark"    orm:"remark"`
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at"`
}
