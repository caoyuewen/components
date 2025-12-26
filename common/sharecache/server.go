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
	"math/rand"
	"time"
)

var ServerMgr = Server{AliveTime: 0}

type Server struct {
	AliveTime int64 // LastCheckTime 15 以内算存活
}

type ServerInfo struct {
	Id            string // 唯一 serverId
	ServerName    string // 服务名称
	Domain        string // 域名
	Status        int    // 状态 0 可用 1 链接占满 2 shutdown
	LastCheckTime int64  // 上次检查状态时间
}

func (t *Server) getServerKey() string {
	return "Server"
}

func (t *Server) GetServerList() ([]ServerInfo, error) {
	var res []ServerInfo
	key := t.getServerKey()

	// 从 Redis 哈希中获取所有服务器信息
	serverMap, err := dbredis.Client().HGetAll(context.Background(), key).Result()
	if err != nil {
		log.Errorf("cache GetServerList HGetAll err:%v", err)
		return res, err
	}

	// 将哈希字段转换为 ServerInfo 切片
	for _, v := range serverMap {
		var serverInfo ServerInfo
		err := json.Unmarshal([]byte(v), &serverInfo)
		if err != nil {
			log.Errorf("cache GetServerList unmarshal err:%v ; result:%s", err, v)
			return res, err
		}
		res = append(res, serverInfo)
	}
	return res, nil
}

func (t *Server) GetServerById(id string) (ServerInfo, error) {
	var res ServerInfo
	key := t.getServerKey()

	// 从 Redis 哈希中获取特定服务器信息
	result, err := dbredis.Client().HGet(context.Background(), key, id).Result()
	if err != nil {
		log.Infof("cache GetServerById HGet err:%v ; gen:%s", err, id)
		return res, err
	}

	err = json.Unmarshal([]byte(result), &res)
	if err != nil {
		log.Errorf("cache GetServerById unmarshal err:%v ; result:%s", err, result)
		return res, err
	}
	return res, nil
}

func (t *Server) RegisterServer(info ServerInfo) error {
	gs, err := t.GetServerById(info.Id)
	if errors.Is(err, redis.Nil) {
		return t.setServerInfo(info)
	} else if err != nil {
		return err
	}

	if time.Now().Unix()-gs.LastCheckTime <= t.AliveTime {
		return fmt.Errorf("game server still alive %s", util.ToJson(gs))
	}

	return t.setServerInfo(info)
}

func (t *Server) CheckingOrUpdate(info ServerInfo) error {
	return t.setServerInfo(info)
}

func (t *Server) setServerInfo(info ServerInfo) error {
	key := t.getServerKey()

	// 将 ServerInfo 对象编码为 JSON
	data, err := json.Marshal(info)
	if err != nil {
		log.Errorf("cache SetServerInfo marshal err:%v", err)
		return err
	}

	// 将服务器信息存储到 Redis 哈希中
	_, err = dbredis.Client().HSet(context.Background(), key, info.Id, data).Result()
	if err != nil {
		log.Errorf("cache SetServerInfo HSet err:%v", err)
		return err
	}
	return nil
}

func (t *Server) GetRandomAvailableServer() (ServerInfo, error) {
	var availableServers []ServerInfo
	key := t.getServerKey()

	// 从 Redis 哈希中获取所有服务器信息
	serverMap, err := dbredis.Client().HGetAll(context.Background(), key).Result()
	if err != nil {
		log.Printf("cache GetRandomAvailableServer HGetAll err:%v", err)
		return ServerInfo{}, err
	}

	// 检查是否返回了空的服务器列表
	if len(serverMap) == 0 {
		return ServerInfo{}, fmt.Errorf("no servers found in the cache")
	}

	currentTime := time.Now().Unix()

	// 遍历所有服务器，寻找状态为可用且 LastCheckTime 在 15 秒以内的服务器
	for _, v := range serverMap {
		var serverInfo ServerInfo
		err := json.Unmarshal([]byte(v), &serverInfo)
		if err != nil {
			log.Printf("cache GetRandomAvailableServer unmarshal err:%v ; result:%s", err, v)
			continue
		}

		// 检查服务器状态是否为可用，且 LastCheckTime 在规定时间以内
		if serverInfo.Status == 0 && currentTime-serverInfo.LastCheckTime <= t.AliveTime {
			availableServers = append(availableServers, serverInfo)
		}
	}

	// 检查是否找到了可用的服务器
	if len(availableServers) == 0 {
		return ServerInfo{}, fmt.Errorf("no available bullet server found")
	}

	// 从可用服务器列表中随机选择一个
	rand.NewSource(time.Now().UnixNano()) // 设定随机种子
	selectedServer := availableServers[rand.Intn(len(availableServers))]

	return selectedServer, nil
}
