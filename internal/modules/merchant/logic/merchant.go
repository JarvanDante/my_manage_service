// Package logic 商户模块业务实现(M1)。
package logic

import (
	"context"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/gtime"

	"github.com/JarvanDante/my_manage_service/internal/model/entity"
	"github.com/JarvanDante/my_manage_service/internal/modules/merchant/domain"
	"github.com/JarvanDante/my_manage_service/internal/modules/merchant/service"
)

type sMerchant struct {
	repo domain.Repository
}

func New(repo domain.Repository) service.IMerchant { return &sMerchant{repo: repo} }

func (s *sMerchant) List(ctx context.Context, in service.ListInput) (*service.ListDTO, error) {
	if in.Page <= 0 {
		in.Page = 1
	}
	if in.Size <= 0 {
		in.Size = 20
	}
	list, total, err := s.repo.List(ctx, domain.Filter{Keyword: in.Keyword, Status: in.Status}, in.Page, in.Size)
	if err != nil {
		return nil, err
	}
	out := make([]*service.MerchantDTO, 0, len(list))
	for _, m := range list {
		out = append(out, toDTO(m))
	}
	return &service.ListDTO{List: out, Total: total, Page: in.Page, Size: in.Size}, nil
}

func (s *sMerchant) Create(ctx context.Context, in service.SaveInput) (int64, error) {
	if in.Name == "" {
		return 0, gerror.New("商户名必填")
	}
	exist, err := s.repo.FindByName(ctx, in.Name)
	if err != nil {
		return 0, err
	}
	if exist != nil {
		return 0, gerror.New("商户名已存在")
	}
	return s.repo.Create(ctx, &entity.Merchant{
		Name: in.Name, Contact: in.Contact, Phone: in.Phone, Remark: in.Remark,
	})
}

func (s *sMerchant) Update(ctx context.Context, in service.SaveInput) error {
	if in.Id <= 0 {
		return gerror.New("商户ID无效")
	}
	if in.Name == "" {
		return gerror.New("商户名必填")
	}
	m, err := s.repo.FindById(ctx, in.Id)
	if err != nil {
		return err
	}
	if m == nil {
		return gerror.New("商户不存在")
	}
	// 名称查重(排除自己)
	exist, err := s.repo.FindByName(ctx, in.Name)
	if err != nil {
		return err
	}
	if exist != nil && exist.Id != in.Id {
		return gerror.New("商户名已存在")
	}
	return s.repo.Update(ctx, &entity.Merchant{
		Id: in.Id, Name: in.Name, Contact: in.Contact, Phone: in.Phone, Remark: in.Remark,
	})
}

// SetStatus 启用(1)/停用(0)。不物理删除。
// TODO(M2): 停用前校验名下站点状态, 提示影响范围。
func (s *sMerchant) SetStatus(ctx context.Context, id int64, status int) error {
	if id <= 0 {
		return gerror.New("商户ID无效")
	}
	if status != 0 && status != 1 {
		return gerror.New("status 仅支持 0/1")
	}
	m, err := s.repo.FindById(ctx, id)
	if err != nil {
		return err
	}
	if m == nil {
		return gerror.New("商户不存在")
	}
	return s.repo.SetStatus(ctx, id, status)
}

func toDTO(m *entity.Merchant) *service.MerchantDTO {
	return &service.MerchantDTO{
		Id: m.Id, Name: m.Name, Contact: m.Contact, Phone: m.Phone,
		Status: m.Status, Remark: m.Remark, CreatedAt: fmtTime(m.CreatedAt),
	}
}

func fmtTime(t *gtime.Time) string {
	if t == nil {
		return ""
	}
	return t.String()
}
