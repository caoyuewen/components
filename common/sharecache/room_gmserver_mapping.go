package sharecache

import (
	"context"
	"errors"
	"github.com/caoyuewen/components/dbs/dbredis"
	"github.com/go-redis/redis/v8"
	log "github.com/sirupsen/logrus"
)

var RoomGmServerCacheMgr RoomGmServerCache

type RoomGmServerCache struct {
}

type RoomGmServer struct {
	RoomId       string
	GmHttpDomain string
}

// 获取 Redis 中 Hash 键的名称
func (t *RoomGmServerCache) getMapKey() string {
	return "RoomGmServerTable" // 键名，用于存储房间与话题的映射
}

// SetRoomGmServer 将房间ID与话题存储在 Redis 的 Hash 中
func (t *RoomGmServerCache) SetRoomGmServer(roomId string, topic string) error {
	key := t.getMapKey() // Hash 键名

	// 将 roomId 和 topic 作为键值对存储在 Redis Hash 中
	result := dbredis.Client().HSet(context.Background(), key, roomId, topic)
	if result.Err() != nil {
		log.Errorf("cache SetRoomGmServer err:%v ; roomId:%s, topic:%s", result.Err(), roomId, topic)
		return result.Err()
	}

	// 可选：设置整个 Hash 的过期时间
	//expire := 24 * time.Hour
	//dbredis.Client().Expire(context.Background(), key, expire)

	return nil
}

// GetGmServerByRoomId 根据 roomId 从 Redis Hash 中获取对应的 topic
func (t *RoomGmServerCache) GetGmServerByRoomId(roomId string) (string, error) {
	key := t.getMapKey() // Hash 键名

	// 根据 roomId 从 Hash 中获取对应的 topic
	result, err := dbredis.Client().HGet(context.Background(), key, roomId).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			log.Warnf("cache GetRoomGmServer: roomId %s not found", roomId)
			return "", nil
		}
		log.Errorf("cache GetRoomGmServer err:%v ; roomId:%s", err, roomId)
		return "", err
	}

	return result, nil
}

// DeleteRoomGmServer 根据 roomId 删除对应的 topic
func (t *RoomGmServerCache) DeleteRoomGmServer(roomId string) error {
	key := t.getMapKey() // Hash 键名

	// 删除 Hash 中指定的 roomId 键值对
	result := dbredis.Client().HDel(context.Background(), key, roomId)
	if result.Err() != nil {
		log.Errorf("cache DeleteRoomGmServer err:%v ; roomId:%s", result.Err(), roomId)
		return result.Err()
	}

	return nil
}
