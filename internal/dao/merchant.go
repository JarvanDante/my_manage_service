package dao

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	"github.com/JarvanDante/my_manage_service/internal/model/entity"
	merchdomain "github.com/JarvanDante/my_manage_service/internal/modules/merchant/domain"
)

type merchantRepo struct{}

// NewMerchantRepo 返回 merchant 领域仓储实现。
func NewMerchantRepo() merchdomain.Repository { return &merchantRepo{} }

func (r *merchantRepo) List(ctx context.Context, f merchdomain.Filter, page, size int) ([]*entity.Merchant, int, error) {
	m := g.Model("merchant").Ctx(ctx)
	if f.Keyword != "" {
		kw := "%" + f.Keyword + "%"
		m = m.Where("(name ILIKE ? OR contact ILIKE ?)", kw, kw)
	}
	switch f.Status {
	case 1:
		m = m.Where("status", 1)
	case 2:
		m = m.Where("status", 0)
	}
	total, err := m.Clone().Count()
	if err != nil {
		return nil, 0, err
	}
	var list []*entity.Merchant
	err = m.Clone().OrderDesc("id").Page(page, size).Scan(&list)
	return list, total, err
}

func (r *merchantRepo) FindById(ctx context.Context, id int64) (*entity.Merchant, error) {
	var m *entity.Merchant
	err := g.Model("merchant").Ctx(ctx).Where("id", id).Scan(&m)
	return m, err
}

func (r *merchantRepo) FindByName(ctx context.Context, name string) (*entity.Merchant, error) {
	var m *entity.Merchant
	err := g.Model("merchant").Ctx(ctx).Where("name", name).Scan(&m)
	return m, err
}

func (r *merchantRepo) Create(ctx context.Context, e *entity.Merchant) (int64, error) {
	res, err := g.Model("merchant").Ctx(ctx).Data(g.Map{
		"name": e.Name, "contact": e.Contact, "phone": e.Phone,
		"status": 1, "remark": e.Remark,
	}).Insert()
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *merchantRepo) Update(ctx context.Context, e *entity.Merchant) error {
	_, err := g.Model("merchant").Ctx(ctx).Where("id", e.Id).Data(g.Map{
		"name": e.Name, "contact": e.Contact, "phone": e.Phone,
		"remark": e.Remark, "updated_at": gtime.Now(),
	}).Update()
	return err
}

func (r *merchantRepo) SetStatus(ctx context.Context, id int64, status int) error {
	_, err := g.Model("merchant").Ctx(ctx).Where("id", id).Data(g.Map{
		"status": status, "updated_at": gtime.Now(),
	}).Update()
	return err
}

// CountActiveSites 商户名下未停用的站点数(控制面同库, dao 层内联查询)。
func (r *merchantRepo) CountActiveSites(ctx context.Context, merchantId int64) (int, error) {
	return g.Model("site").Ctx(ctx).Where("merchant_id", merchantId).WhereNot("status", 2).Count()
}
