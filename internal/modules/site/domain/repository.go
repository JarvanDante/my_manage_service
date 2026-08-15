// Package domain 站点领域层(M2)。
package domain

import (
	"context"

	"github.com/JarvanDante/my_manage_service/internal/model/entity"
)

// Filter 站点列表筛选。
type Filter struct {
	Keyword    string // site_code/名称 模糊
	MerchantId int64
	Env        string
	Status     int // -1全部 0筹备 1上线 2停用(注意0是有效状态)
}

// SiteItem 列表项(联表商户名)。
type SiteItem struct {
	entity.Site
	MerchantName string `orm:"merchant_name"`
}

type Repository interface {
	List(ctx context.Context, f Filter, page, size int) ([]*SiteItem, int, error)
	FindById(ctx context.Context, id int64) (*entity.Site, error)
	FindItemById(ctx context.Context, id int64) (*SiteItem, error)
	FindByCode(ctx context.Context, code string) (*entity.Site, error)
	Create(ctx context.Context, s *entity.Site) (int64, error)
	Update(ctx context.Context, s *entity.Site) error
	UpdateAppCredentials(ctx context.Context, id int64, appKey, encSecret string) error
	SetStatus(ctx context.Context, id int64, status int) error
	MerchantExists(ctx context.Context, merchantId int64) (bool, error)

	// 域名
	Domains(ctx context.Context, siteId int64) ([]*entity.SiteDomain, error)
	FindDomain(ctx context.Context, domain string) (*entity.SiteDomain, error)
	FindDomainById(ctx context.Context, id int64) (*entity.SiteDomain, error)
	BindDomain(ctx context.Context, siteId int64, domain, kind string, https, isMain int) (int64, error)
	UnbindDomain(ctx context.Context, id int64) error
	SetMainDomain(ctx context.Context, siteId, domainId int64) error
}
