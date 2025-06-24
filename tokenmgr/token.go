package tokenmgr

import (
	"encoding/json"
	"github.com/caoyuewen/components/util"
	log "github.com/sirupsen/logrus"
	"time"
)

type TokenInfo struct {
	AppName  string
	UserId   string
	ExpireAt int64
}

func CreateToken(aes string, info TokenInfo) (string, error) {
	m := TokenInfo{
		UserId:   info.UserId,
		AppName:  info.AppName,
		ExpireAt: info.ExpireAt,
	}
	data, err := json.Marshal(m)
	if err != nil {
		log.Info("create token err", err)
		return "", err
	}
	encrypted, err := util.AESEncrypt(string(data), aes)
	if err != nil {
		log.Info("create token err", err)
		return "", err
	}
	return encrypted, nil
}

func DecryptToken(aes, token string) (*TokenInfo, error) {
	bytes, err := util.AESDecrypt(token, aes)
	if err != nil {
		log.Info("decrypt token err", err)
		return nil, err
	}

	var info TokenInfo
	err = json.Unmarshal([]byte(bytes), &info)
	if err != nil {
		log.Info("decrypt token unmarshal err", err)
		return nil, err
	}
	return &info, nil
}

func VerifyToken(aes, token, appName, uid string, loc *time.Location) bool {
	info, err := DecryptToken(aes, token)
	if err != nil {
		return false
	}

	if info == nil {
		return false
	}

	if info.ExpireAt < time.Now().In(loc).Unix() { // 过期
		return false
	}

	if info.AppName != appName {
		return false
	}

	if info.UserId != uid {
		return false
	}

	return true
}

func (t *TokenInfo) IsExpired(loc *time.Location) bool {
	return time.Now().In(loc).Unix() >= t.ExpireAt
}
