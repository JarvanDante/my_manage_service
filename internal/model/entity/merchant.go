// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
package entity

import "github.com/gogf/gf/v2/os/gtime"

type Merchant struct {
	Id        int64       `json:"id"        orm:"id"`
	Name      string      `json:"name"      orm:"name"`
	Contact   string      `json:"contact"   orm:"contact"`
	Phone     string      `json:"phone"     orm:"phone"`
	Status    int         `json:"status"    orm:"status"`
	Remark    string      `json:"remark"    orm:"remark"`
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at"`
	UpdatedAt *gtime.Time `json:"updatedAt" orm:"updated_at"`
}
