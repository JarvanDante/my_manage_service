// Package logic 总后台管理员业务实现。
package logic

import (
	"context"
	"strings"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/gtime"

	"github.com/JarvanDante/my_manage_service/internal/model/entity"
	"github.com/JarvanDante/my_manage_service/internal/modules/auth/domain"
	"github.com/JarvanDante/my_manage_service/internal/modules/auth/service"
	"github.com/JarvanDante/my_manage_service/internal/shared/kit"
)

// issueToken 签发会话; 单测可替换, 避免依赖 Redis。
var issueToken = kit.IssueToken

type sAuth struct {
	repo    domain.Repository
	pending totpPending
}

func New(repo domain.Repository) service.IAuth {
	return &sAuth{repo: repo, pending: redisTotpPending{}}
}

func newAuth(repo domain.Repository, pending totpPending) *sAuth {
	return &sAuth{repo: repo, pending: pending}
}

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

	bound := totpBound(a)
	code := strings.TrimSpace(in.TotpCode)
	if code == "" {
		return s.challengeTotp(ctx, a, bound)
	}
	if err = s.verifyTotp(ctx, a, bound, code); err != nil {
		return nil, err
	}

	_ = s.repo.UpdateLoginInfo(ctx, a.Id, in.Ip)
	token, err := issueToken(ctx, a.Id)
	if err != nil {
		return nil, err
	}
	return &service.LoginDTO{Token: token, Admin: toInfo(a)}, nil
}

func (s *sAuth) challengeTotp(ctx context.Context, a *entity.ManageAdmin, bound bool) (*service.LoginDTO, error) {
	if bound {
		return &service.LoginDTO{NeedTotp: true, TotpBound: true}, nil
	}
	secret, err := s.pending.Get(ctx, a.Id)
	if err != nil {
		return nil, err
	}
	if secret == "" {
		secret, _, _, err = kit.GenerateTOTP(a.Username)
		if err != nil {
			return nil, err
		}
		if err = s.pending.Set(ctx, a.Id, secret); err != nil {
			return nil, err
		}
	}
	_, qr, err := kit.TOTPQRFromSecret(a.Username, secret)
	if err != nil {
		return nil, err
	}
	return &service.LoginDTO{
		NeedTotp:   true,
		TotpBound:  false,
		TotpQR:     qr,
		TotpSecret: secret,
	}, nil
}

func (s *sAuth) verifyTotp(ctx context.Context, a *entity.ManageAdmin, bound bool, code string) error {
	if bound {
		if !kit.ValidateTOTP(code, a.TotpSecret) {
			return gerror.New("验证码错误")
		}
		return nil
	}
	secret, err := s.pending.Get(ctx, a.Id)
	if err != nil {
		return err
	}
	if secret == "" {
		return gerror.New("请先扫描二维码绑定谷歌验证器")
	}
	if !kit.ValidateTOTP(code, secret) {
		return gerror.New("验证码错误")
	}
	if err = s.repo.BindTotp(ctx, a.Id, secret); err != nil {
		return err
	}
	_ = s.pending.Del(ctx, a.Id)
	a.TotpSecret = secret
	a.TotpBoundAt = gtime.Now()
	return nil
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

func (s *sAuth) ResetTotp(ctx context.Context, id int64) error {
	if id <= 0 {
		return gerror.New("管理员ID无效")
	}
	a, err := s.repo.FindById(ctx, id)
	if err != nil {
		return err
	}
	if a == nil {
		return gerror.New("管理员不存在")
	}
	if err = s.repo.ResetTotp(ctx, id); err != nil {
		return err
	}
	_ = s.pending.Del(ctx, id)
	return nil
}

func totpBound(a *entity.ManageAdmin) bool {
	return a != nil && a.TotpSecret != "" && a.TotpBoundAt != nil
}

func toInfo(a *entity.ManageAdmin) *service.AdminInfoDTO {
	return &service.AdminInfoDTO{Id: a.Id, Username: a.Username, Nickname: a.Nickname, RoleId: a.RoleId}
}
