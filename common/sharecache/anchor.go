package sharecache

import (
	"components/dbs/dbredis"
	"components/util"
	"context"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"time"
)

var AnchorCacheMgr AnchorCache

type AnchorCache struct {
}

type AnchorInfo struct {
	AnchorOpenId string
	LiveToken    string
	RoomId       string
	GameToken    string
}

func (t *AnchorCache) getAnchorKey(anchorOpenId string) string {
	return fmt.Sprintf("Anchor.Info.%s", anchorOpenId)
}

func (t *AnchorCache) SetAnchorLoginToken(anchor AnchorInfo, expire time.Duration) error {
	anchorJson := util.ToJson(anchor)
	key := t.getAnchorKey(anchor.AnchorOpenId)
	result := dbredis.Client().Set(context.Background(), key, anchorJson, expire)
	if result.Err() != nil {
		log.Errorf("cache SetAnchorLoginToken err:%v ; anchor :%s", result.Err(), anchorJson)
		result.Err()
	}
	return nil
}

func (t *AnchorCache) GetAnchorLoginToken(anchorOpenId string) (AnchorInfo, error) {
	var res AnchorInfo
	key := t.getAnchorKey(anchorOpenId)
	result, err := dbredis.Client().Get(context.Background(), key).Result()
	if err != nil {
		log.Errorf("cache GetAnchorLoginToken err:%v ; anchorOpenId:%s", err, anchorOpenId)
		return res, err
	}

	err = json.Unmarshal([]byte(result), &res)
	if err != nil {
		log.Errorf("cache GetAnchorLoginToken unmarshal err:%v ; result:%s", err, result)
		return res, err
	}
	return res, nil
}
