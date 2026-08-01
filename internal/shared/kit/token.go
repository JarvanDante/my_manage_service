// Package kit 通用小工具(token / 加密)。
package kit

import (
	"context"
	"strings"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/gogf/gf/v2/util/guid"
)

const (
	tokenPrefix = "manage:token:" // token   -> adminId
	aidPrefix   = "manage:aid:"   // adminId -> 当前 token(单点登录)
	tokenTTL    = int64(7 * 24 * 3600)
)

func normalize(token string) string {
	return strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(token), "Bearer "))
}

// IssueToken 签发总后台管理员 token(单点: 踢旧)。
func IssueToken(ctx context.Context, adminId int64) (string, error) {
	aidKey := aidPrefix + gconv.String(adminId)
	if old, err := g.Redis().Do(ctx, "GET", aidKey); err == nil && !old.IsNil() && old.String() != "" {
		_, _ = g.Redis().Do(ctx, "DEL", tokenPrefix+old.String())
	}
	token := guid.S()
	if _, err := g.Redis().Do(ctx, "SETEX", tokenPrefix+token, tokenTTL, adminId); err != nil {
		return "", err
	}
	if _, err := g.Redis().Do(ctx, "SETEX", aidKey, tokenTTL, token); err != nil {
		return "", err
	}
	return token, nil
}

// ParseToken 校验 token 返回 adminId, 滑动续期。
func ParseToken(ctx context.Context, token string) (int64, error) {
	token = normalize(token)
	if token == "" {
		return 0, gerror.New("empty token")
	}
	key := tokenPrefix + token
	v, err := g.Redis().Do(ctx, "GET", key)
	if err != nil {
		return 0, err
	}
	if v.IsNil() || v.String() == "" {
		return 0, gerror.New("token not found or expired")
	}
	id := v.Int64()
	_, _ = g.Redis().Do(ctx, "EXPIRE", key, tokenTTL)
	_, _ = g.Redis().Do(ctx, "EXPIRE", aidPrefix+gconv.String(id), tokenTTL)
	return id, nil
}

// RevokeByAdminId 按管理员ID撤销当前 token(退出登录)。
func RevokeByAdminId(ctx context.Context, adminId int64) error {
	aidKey := aidPrefix + gconv.String(adminId)
	if v, err := g.Redis().Do(ctx, "GET", aidKey); err == nil && !v.IsNil() && v.String() != "" {
		_, _ = g.Redis().Do(ctx, "DEL", tokenPrefix+v.String())
	}
	_, err := g.Redis().Do(ctx, "DEL", aidKey)
	return err
}
