package sharecache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/caoyuewen/components/dbs/dbredis"
	"github.com/caoyuewen/components/util"
	"github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"
	"math"
	"time"
)

var GameServerMgr = GameServer{AliveTime: 15}

type GameServer struct {
	AliveTime int64 // LastCheckTime 15 以内算存活
}

type GameServerInfo struct {
	Id             string // 唯一 serverId
	GameServerName string // 游戏服务名称
	CurrentConnNum int    // 当前连接数
	Domain         string // ws域名链接地址
	Status         int    // 状态 0 可用 1 链接占满 2 shutdown
	LastCheckTime  int64  // 上次检查状态时间
}

func (t *GameServer) getServerKey() string {
	return "GameServer"
}

func (t *GameServer) GetGameServerList() ([]GameServerInfo, error) {
	var res []GameServerInfo
	key := t.getServerKey()

	// 从 Redis 哈希中获取所有服务器信息
	serverMap, err := dbredis.Client().HGetAll(context.Background(), key).Result()
	if err != nil {
		log.Errorf("cache GetGameServerList HGetAll err:%v", err)
		return res, err
	}

	// 将哈希字段转换为 GameServerInfo 切片
	for _, v := range serverMap {
		var serverInfo GameServerInfo
		err := json.Unmarshal([]byte(v), &serverInfo)
		if err != nil {
			log.Errorf("cache GetGameServerList unmarshal err:%v ; result:%s", err, v)
			return res, err
		}
		res = append(res, serverInfo)
	}
	return res, nil
}

func (t *GameServer) GetGameServerById(id string) (GameServerInfo, error) {
	var res GameServerInfo
	key := t.getServerKey()

	// 从 Redis 哈希中获取特定服务器信息
	result, err := dbredis.Client().HGet(context.Background(), key, id).Result()
	if err != nil {
		log.Infof("cache GetGameServerById HGet err:%v ; gen:%s", err, id)
		return res, err
	}

	err = json.Unmarshal([]byte(result), &res)
	if err != nil {
		log.Errorf("cache GetGameServerById unmarshal err:%v ; result:%s", err, result)
		return res, err
	}
	return res, nil
}

func (t *GameServer) RegisterGameServer(info GameServerInfo) error {
	gs, err := t.GetGameServerById(info.Id)
	if errors.Is(err, redis.Nil) {
		return t.setGameServerInfo(info)
	} else if err != nil {
		return err
	}

	if time.Now().Unix()-gs.LastCheckTime <= t.AliveTime {
		return fmt.Errorf("game server still alive %s", util.ToJson(gs))
	}

	return t.setGameServerInfo(info)
}

func (t *GameServer) CheckingOrUpdate(info GameServerInfo) error {
	return t.setGameServerInfo(info)
}

func (t *GameServer) setGameServerInfo(info GameServerInfo) error {
	key := t.getServerKey()

	// 将 GameServerInfo 对象编码为 JSON
	data, err := json.Marshal(info)
	if err != nil {
		log.Errorf("cache SetGameServerInfo marshal err:%v", err)
		return err
	}

	// 将服务器信息存储到 Redis 哈希中
	_, err = dbredis.Client().HSet(context.Background(), key, info.Id, data).Result()
	if err != nil {
		log.Errorf("cache SetGameServerInfo HSet err:%v", err)
		return err
	}
	return nil
}

func (t *GameServer) GetLeastConnAvailableGameServer() (GameServerInfo, error) {
	var leastConnServer GameServerInfo
	key := t.getServerKey()

	// 从 Redis 哈希中获取所有服务器信息
	serverMap, err := dbredis.Client().HGetAll(context.Background(), key).Result()
	if err != nil {
		log.Printf("cache GetLeastConnAvailableGameServer HGetAll err:%v", err)
		return leastConnServer, err
	}

	// 检查是否返回了空的服务器列表
	if len(serverMap) == 0 {
		return leastConnServer, fmt.Errorf("no servers found in the cache")
	}

	// 初始化最小连接数为最大值
	minConn := math.MaxInt32
	found := false
	currentTime := time.Now().Unix()

	// 遍历所有服务器，寻找连接数最少、状态为可用且LastCheckTime在15秒以内的服务器
	for _, v := range serverMap {
		var serverInfo GameServerInfo
		err := json.Unmarshal([]byte(v), &serverInfo)
		if err != nil {
			log.Printf("cache GetLeastConnAvailableGameServer unmarshal err:%v ; result:%s", err, v)
			continue
		}

		// 检查服务器状态是否为可用，且LastCheckTime在15秒以内
		if serverInfo.Status == 0 && currentTime-serverInfo.LastCheckTime <= t.AliveTime && serverInfo.CurrentConnNum < minConn {
			minConn = serverInfo.CurrentConnNum
			leastConnServer = serverInfo
			found = true
		}
	}

	if !found {
		return leastConnServer, fmt.Errorf("no available game server found")
	}

	return leastConnServer, nil
}
