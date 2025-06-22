package tictok

import (
	"components/dbs/dbredis"
	"components/third/tictok/tictokapi"
	"components/util"
	"context"
	log "github.com/sirupsen/logrus"
	"time"
)

const (
	CheckAccessTokenInterval = 1                   // 单位分钟 每分钟检测一次 如果过期时间在5分钟以内就要更新token
	RedisAccessTokenKey      = "TictokAccessToken" // Redis 中存储 access token 的 key
)

var AccessTokenMgr accessToken

type accessToken struct {
	appId  string
	secret string
}

func InitAccessToken(appid, secret string) {
	AccessTokenMgr.appId = appid
	AccessTokenMgr.secret = secret
	getAccessToken(appid, secret)
	go flushAccessToken(appid, secret)
}

func (a *accessToken) GetAccessToken() (string, error) {
	// 从 Redis 中获取 access token
	token, err := dbredis.Client().Get(context.TODO(), RedisAccessTokenKey).Result()
	if err != nil {
		log.Error("get access token from api,get token in redis err:", err.Error())
		return getAccessToken(a.appId, a.secret), nil
	}
	return token, nil
}

func getAccessToken(appid, secret string) string {
	resp, err := tictokapi.GetAccessToken(appid, secret)
	if err != nil {
		log.Errorf("getAccessToken err:%s", err.Error())
	} else if resp.ErrNo != tictokapi.ApiCodeSuccess {
		log.Errorf("getAccessToken err:%s", util.ToJson(resp))
	}

	// 将 access token 写入 Redis
	err = dbredis.Client().Set(context.TODO(), RedisAccessTokenKey, resp.Data.AccessToken, time.Second*time.Duration(resp.Data.ExpiresIn)).Err()
	if err != nil {
		log.Errorf("failed to save access token to redis: %v", err)
	} else {
		log.WithField("resp", resp).Trace("get tictok access info and saved to redis")
	}

	log.Debug("getAccessToken resp:", util.ToJson(resp))
	return resp.Data.AccessToken
}

func flushAccessToken(appid, secret string) {
	ticker := time.NewTicker(CheckAccessTokenInterval * time.Minute)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			//log.Info("Access token checking...")
			ttl, err := dbredis.Client().TTL(context.TODO(), RedisAccessTokenKey).Result()
			if err != nil {
				log.Error("check access token err:", err)
				continue
			}
			// 检查剩余过期时间是否在 5 分钟以内
			if ttl > 0 && ttl <= 5*time.Minute {
				log.Info("Access token will expire soon, refreshing...")
				getAccessToken(appid, secret)
			} else if ttl < 0 {
				// 当 ttl < 0 时表示 key 不存在，需要立即获取新的 token
				log.Info("Access token does not exist, fetching a new one...")
				getAccessToken(appid, secret)
			}
		}
	}
}
