// Package logic 总后台管理员业务实现。
package logic

import (
	"context"

	"github.com/gogf/gf/v2/errors/gerror"

	"github.com/JarvanDante/my_manage_service/internal/model/entity"
	"github.com/JarvanDante/my_manage_service/internal/modules/auth/domain"
	"github.com/JarvanDante/my_manage_service/internal/modules/auth/service"
	"github.com/JarvanDante/my_manage_service/internal/shared/kit"
)

type sAuth struct {
	repo domain.Repository
}

func New(repo domain.Repository) service.IAuth { return &sAuth{repo: repo} }

func (s *sAuth) Login(ctx context.Context, in service.LoginInput) (*service.LoginDTO, error) {
	if in.Username == "" || in.Password == "" {
		return nil, gerror.New("账号或密码不能为空")
	}
	a, err := s.repo.FindByUsername(ctx, in.Username)
	if err != nil {
		return nil, err
	}
	if a == nil {
		return nil, gerror.New("账号或密码错误")
	}
	if a.Status != 1 {
		return nil, gerror.New("账号已被禁用")
	}
	if !kit.CheckPassword(in.Password, a.Password) {
		return nil, gerror.New("账号或密码错误")
	}
	_ = s.repo.UpdateLoginInfo(ctx, a.Id, in.Ip)
	token, err := kit.IssueToken(ctx, a.Id)
	if err != nil {
		return nil, err
	}
	return &service.LoginDTO{Token: token, Admin: toInfo(a)}, nil
}

func (s *sAuth) Logout(ctx context.Context, adminId int64) error {
	return kit.RevokeByAdminId(ctx, adminId)
}

func (s *sAuth) Info(ctx context.Context, adminId int64) (*service.AdminInfoDTO, error) {
	a, err := s.repo.FindById(ctx, adminId)
	if err != nil {
		return nil, err
	}
	if a == nil {
		return nil, gerror.New("管理员不存在")
	}
	return toInfo(a), nil
}

func toInfo(a *entity.ManageAdmin) *service.AdminInfoDTO {
	return &service.AdminInfoDTO{Id: a.Id, Username: a.Username, Nickname: a.Nickname, RoleId: a.RoleId}
}
