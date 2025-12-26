package dbredis

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"
)

var (
	client        atomic.Pointer[redis.Client]
	clusterClient atomic.Pointer[redis.ClusterClient]
	once          sync.Once
	isCluster     atomic.Bool
	initialized   atomic.Bool
)


// StartUp initializes the global Redis client or cluster client with the provided options and starts the connection checker
func StartUp(addrs []string, checkInterval time.Duration) {
	once.Do(func() {
		if len(addrs) > 1 {
			isCluster.Store(true)
			cc := redis.NewClusterClient(&redis.ClusterOptions{
				Addrs: addrs,
			})

			ctx := context.Background()
			_, err := cc.Ping(ctx).Result()
			if err != nil {
				log.Errorf("[REDIS] Failed to connected redis cluster:%v err:%s ", addrs, err.Error())
				panic(err.Error())
			}

			clusterClient.Store(cc)
			initialized.Store(true)
			log.Info("[REDIS] Redis cluster connected successfully")
			go connectionCheckerCluster(&redis.ClusterOptions{Addrs: addrs}, checkInterval)
		} else {
			isCluster.Store(false)
			c := redis.NewClient(&redis.Options{
				Addr: addrs[0],
			})

			ctx := context.Background()
			_, err := c.Ping(ctx).Result()
			if err != nil {
				log.Error("[REDIS] Failed to connect redis: " + err.Error())
				panic(err.Error())
			}

			client.Store(c)
			initialized.Store(true)
			log.Info("[REDIS] Redis connected successfully:", addrs)
			go connectionChecker(&redis.Options{Addr: addrs[0]}, checkInterval)
		}
	})
}

func connectionChecker(options *redis.Options, checkInterval time.Duration) {
	ticker := time.NewTicker(checkInterval)
	defer ticker.Stop()

	for range ticker.C {
		ctx := context.Background()
		c := client.Load()
		if c == nil {
			continue
		}
		_, err := c.Ping(ctx).Result()
		if err != nil {
			log.Warn("[REDIS] Redis connection lost, attempting to reconnect...")
			reconnect(options)
		}
	}
}

func connectionCheckerCluster(options *redis.ClusterOptions, checkInterval time.Duration) {
	ticker := time.NewTicker(checkInterval)
	defer ticker.Stop()

	for range ticker.C {
		ctx := context.Background()
		cc := clusterClient.Load()
		if cc == nil {
			continue
		}
		_, err := cc.Ping(ctx).Result()
		if err != nil {
			log.Warn("[REDIS] Redis cluster connection lost, attempting to reconnect...")
			reconnectCluster(options)
		}
	}
}

func reconnect(options *redis.Options) {
	for {
		ctx := context.Background()
		newClient := redis.NewClient(options)
		_, err := newClient.Ping(ctx).Result()
		if err == nil {
			client.Store(newClient)
			log.Info("[REDIS] Successfully reconnected to redis")
			break
		}
		log.Error("[REDIS] Reconnect failed, retrying... " + err.Error())
		time.Sleep(5 * time.Second)
	}
}

func reconnectCluster(options *redis.ClusterOptions) {
	for {
		ctx := context.Background()
		newClient := redis.NewClusterClient(options)
		_, err := newClient.Ping(ctx).Result()
		if err == nil {
			clusterClient.Store(newClient)
			log.Info("[REDIS] Successfully reconnected to redis cluster")
			break
		}
		log.Error("[REDIS] Reconnect failed, retrying... " + err.Error())
		time.Sleep(5 * time.Second)
	}
}

// Client 返回统一的 Redis 客户端接口（自动判断单机/集群）
// 如果未初始化会 panic
func Client() redis.UniversalClient {
	if !initialized.Load() {
		panic("[REDIS] Client not initialized, call StartUp first")
	}
	if isCluster.Load() {
		return clusterClient.Load()
	}
	return client.Load()
}

// RawClient 返回单机模式的原始客户端
func RawClient() *redis.Client {
	return client.Load()
}

// RawClusterClient 返回集群模式的原始客户端
func RawClusterClient() *redis.ClusterClient {
	return clusterClient.Load()
}

// IsCluster 返回是否为集群模式
func IsCluster() bool {
	return isCluster.Load()
}

// IsInitialized 返回是否已初始化
func IsInitialized() bool {
	return initialized.Load()
}
