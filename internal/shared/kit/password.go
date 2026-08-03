package kit

import (
	"github.com/gogf/gf/v2/errors/gerror"
	"golang.org/x/crypto/bcrypt"
)

// HashPassword 使用 bcrypt 哈希明文密码(盐内嵌于哈希结果, 无需单独 salt 字段)。
func HashPassword(plain string) (string, error) {
	if plain == "" {
		return "", gerror.New("密码不能为空")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
	if err != nil {
		return "", gerror.Wrap(err, "密码加密失败")
	}
	return string(hash), nil
}

// CheckPassword 校验明文密码与 bcrypt 哈希是否匹配。
func CheckPassword(plain, hash string) bool {
	if plain == "" || hash == "" {
		return false
	}
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(plain)) == nil
}
