package tokenmgr

import (
	"encoding/base64"
	"encoding/json"
	"github.com/caoyuewen/components/util"
	log "github.com/sirupsen/logrus"
	"time"
)

const (
	PermissionsNormal = 1 // 普通用户
	PermissionsTest   = 2 // 测试账号
	PermissionsAdmin  = 3 // 可操作后台
)

type TokenInfo struct {
	Id          string
	AppName     string
	Permissions int
	ExpireAt    int64
}

func CreateToken(aes string, info TokenInfo) (string, error) {
	data, err := json.Marshal(info)
	if err != nil {
		log.Info("create token err", err)
		return "", err
	}
	encrypted, err := util.AESEncrypt(string(data), aes)
	if err != nil {
		log.Info("create token err", err)
		return "", err
	}

	token := base64.URLEncoding.EncodeToString([]byte(encrypted))
	return token, nil
}

func DecryptToken(aes, token string) (TokenInfo, error) {
	// ✅ Cookie 里的 token 是 URL-safe base64，需要先解码
	raw, err := base64.URLEncoding.DecodeString(token)
	if err != nil {
		log.Info("base64 decode token err", err)
		return TokenInfo{}, err
	}

	bytes, err := util.AESDecrypt(string(raw), aes)
	if err != nil {
		log.Info("decrypt token err", err)
		return TokenInfo{}, err
	}

	var info TokenInfo
	err = json.Unmarshal([]byte(bytes), &info)
	if err != nil {
		log.Info("unmarshal token err", err)
		return TokenInfo{}, err
	}
	return info, nil
}

func VerifyToken(aes, token, appName string, loc *time.Location) (TokenInfo, bool) {
	info, err := DecryptToken(aes, token)
	if err != nil {
		return info, false
	}

	if info.ExpireAt < time.Now().In(loc).Unix() { // 过期
		return info, false
	}

	if info.AppName != appName {
		return info, false
	}

	return info, true
}

func (t *TokenInfo) IsExpired(loc *time.Location) bool {
	return time.Now().In(loc).Unix() >= t.ExpireAt
}
