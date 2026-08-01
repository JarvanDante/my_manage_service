// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
package entity

import "github.com/gogf/gf/v2/os/gtime"

type Site struct {
	Id         int64       `json:"id"         orm:"id"`
	MerchantId int64       `json:"merchantId" orm:"merchant_id"`
	SiteCode   string      `json:"siteCode"   orm:"site_code"`
	Name       string      `json:"name"       orm:"name"`
	Env        string      `json:"env"        orm:"env"`
	Status     int         `json:"status"     orm:"status"`
	DbHost     string      `json:"dbHost"     orm:"db_host"`
	DbPort     int         `json:"dbPort"     orm:"db_port"`
	DbName     string      `json:"dbName"     orm:"db_name"`
	DbUser     string      `json:"dbUser"     orm:"db_user"`
	DbPass     string      `json:"dbPass"     orm:"db_pass"`
	Remark     string      `json:"remark"     orm:"remark"`
	CreatedAt  *gtime.Time `json:"createdAt"  orm:"created_at"`
	UpdatedAt  *gtime.Time `json:"updatedAt"  orm:"updated_at"`
}
