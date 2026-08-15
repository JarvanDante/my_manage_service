// Package service 站点模块对外接口(M2)。
package service

import "context"

type SiteDTO struct {
	Id           int64
	MerchantId   int64
	MerchantName string
	SiteCode     string
	Name         string
	Env          string
	Status       int
	DbHost       string
	DbPort       int
	DbName       string
	DbUser       string
	DbPassMask   string // 永远脱敏
	AppKey       string // 可展示; secret 永不出现在列表/详情
	Remark       string
	CreatedAt    string
}

// CreateResult 创建站点结果(含明文 APPSECRET, 之后仍可通过 RevealSecret 查看)。
type CreateResult struct {
	Id        int64
	AppKey    string
	AppSecret string
}

// AppCredentialsDTO APPKEY + 明文 APPSECRET。
type AppCredentialsDTO struct {
	AppKey    string
	AppSecret string
}

type DomainDTO struct {
	Id     int64
	Domain string
	Kind   string // admin / h5 / app
	IsMain int
	Https  int
	Status int
}

type SiteDetailDTO struct {
	SiteDTO
	Domains []*DomainDTO
}

type ListInput struct {
	Keyword    string
	MerchantId int64
	Env        string
	Status     int // -1全部 0筹备 1上线 2停用
	Page       int
	Size       int
}

type ListDTO struct {
	List  []*SiteDTO
	Total int
	Page  int
	Size  int
}

type SaveInput struct {
	Id         int64 // 更新时用
	MerchantId int64
	SiteCode   string // 创建后不可改
	Name       string
	Env        string
	DbHost     string
	DbPort     int
	DbName     string
	DbUser     string
	DbPass     string // 明文入参, 落库前 AES; 更新时空=不改
	Remark     string
}

type ProvisionItemDTO struct {
	Key    string // merchant / domain / db / migration / config
	Name   string
	Ok     bool
	Detail string
}

type ProvisionDTO struct {
	AllPassed bool
	Items     []*ProvisionItemDTO
}

type DbCheckDTO struct {
	Ok        bool
	LatencyMs int64
	Message   string
}

type ISite interface {
	List(ctx context.Context, in ListInput) (*ListDTO, error)
	Detail(ctx context.Context, id int64) (*SiteDetailDTO, error)
	Create(ctx context.Context, in SaveInput) (*CreateResult, error)
	Update(ctx context.Context, in SaveInput) error
	SetStatus(ctx context.Context, id int64, status int) error
	DbCheck(ctx context.Context, id int64) (*DbCheckDTO, error)
	ProvisionCheck(ctx context.Context, id int64) (*ProvisionDTO, error)
	RevealSecret(ctx context.Context, id int64) (*AppCredentialsDTO, error)
	ResetSecret(ctx context.Context, id int64) (*AppCredentialsDTO, error)

	BindDomain(ctx context.Context, siteId int64, domain, kind string, https int) (int64, error)
	UnbindDomain(ctx context.Context, siteId, domainId int64) error
	SetMainDomain(ctx context.Context, siteId, domainId int64) error
}
