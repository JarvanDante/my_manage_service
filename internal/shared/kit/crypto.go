package kit

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"io"
	"os"

	"github.com/gogf/gf/v2/errors/gerror"
)

// aesKey 由环境变量 MANAGE_SECRET 派生(sha256 -> 32字节)。
// 用途: 站点数据库密码等敏感登记信息的落库加密(M2 起使用)。
func aesKey() ([]byte, error) {
	secret := os.Getenv("MANAGE_SECRET")
	if secret == "" {
		return nil, gerror.New("环境变量 MANAGE_SECRET 未设置(AES 密钥来源)")
	}
	sum := sha256.Sum256([]byte(secret))
	return sum[:], nil
}

// Encrypt AES-GCM 加密, 输出 base64(nonce+密文)。
func Encrypt(plain string) (string, error) {
	key, err := aesKey()
	if err != nil {
		return "", err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}
	out := gcm.Seal(nonce, nonce, []byte(plain), nil)
	return base64.StdEncoding.EncodeToString(out), nil
}

// Decrypt 解密 Encrypt 的输出。
func Decrypt(enc string) (string, error) {
	key, err := aesKey()
	if err != nil {
		return "", err
	}
	raw, err := base64.StdEncoding.DecodeString(enc)
	if err != nil {
		return "", err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	if len(raw) < gcm.NonceSize() {
		return "", gerror.New("密文格式不正确")
	}
	nonce, data := raw[:gcm.NonceSize()], raw[gcm.NonceSize():]
	plain, err := gcm.Open(nil, nonce, data, nil)
	if err != nil {
		return "", err
	}
	return string(plain), nil
}
