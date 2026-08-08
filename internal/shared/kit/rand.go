package kit

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

// RandomHex 生成 n 字节的十六进制字符串(长度 2n)。
func RandomHex(n int) (string, error) {
	if n <= 0 {
		return "", fmt.Errorf("n must be positive")
	}
	buf := make([]byte, n)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return hex.EncodeToString(buf), nil
}

// NewAppKey 生成 PaaS APPKEY, 形如 ak_ + 24 hex。
func NewAppKey() (string, error) {
	h, err := RandomHex(12)
	if err != nil {
		return "", err
	}
	return "ak_" + h, nil
}

// NewAppSecret 生成 PaaS APPSECRET(32 字节 / 64 hex)。
func NewAppSecret() (string, error) {
	return RandomHex(32)
}
