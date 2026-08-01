// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
package entity

import "github.com/gogf/gf/v2/os/gtime"

// ManageLogBase 总后台操作日志。
type ManageLogBase struct {
	Id        int64       `json:"id"        orm:"id"`
	AdminId   int64       `json:"adminId"   orm:"admin_id"`
	Method    string      `json:"method"    orm:"method"`
	Path      string      `json:"path"      orm:"path"`
	Ip        string      `json:"ip"        orm:"ip"`
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at"`
}
