package dbredis

import (
	"context"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	log "github.com/sirupsen/logrus"
)

var (
	client        *redis.Client
	clusterClient *redis.ClusterClient
	once          sync.Once
	mu            sync.Mutex
	isCluster     bool
)

// StartUp initializes the global Redis client or cluster client with the provided options and starts the connection checker
func StartUp(addrs []string, checkInterval time.Duration) {
	once.Do(func() {
		if len(addrs) > 1 {
			isCluster = true
			clusterClient = redis.NewClusterClient(&redis.ClusterOptions{
				Addrs: addrs,
			})

			ctx := context.Background()
			_, err := clusterClient.Ping(ctx).Result()
			if err != nil {
				log.Errorf("[REDIS] Failed to connect redis cluster: %v err: %s", addrs, err.Error())
				panic(err.Error())
			}

			log.Info("[REDIS] Redis cluster connected successfully")
			go connectionCheckerCluster(&redis.ClusterOptions{Addrs: addrs}, checkInterval)
		} else {
			isCluster = false
			client = redis.NewClient(&redis.Options{
				Addr: addrs[0],
			})

			ctx := context.Background()
			_, err := client.Ping(ctx).Result()
			if err != nil {
				log.Error("[REDIS] Failed to connect redis: " + err.Error())
				return
			}

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
		_, err := client.Ping(ctx).Result()
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
		_, err := clusterClient.Ping(ctx).Result()
		if err != nil {
			log.Warn("[REDIS] Redis cluster connection lost, attempting to reconnect...")
			reconnectCluster(options)
		}
	}
}

func reconnect(options *redis.Options) {
	mu.Lock()
	defer mu.Unlock()

	for {
		ctx := context.Background()
		newClient := redis.NewClient(options)
		_, err := newClient.Ping(ctx).Result()
		if err == nil {
			client = newClient
			log.Info("[REDIS] Successfully reconnected to redis")
			break
		}
		log.Error("[REDIS] Reconnect failed, retrying... " + err.Error())
		time.Sleep(5 * time.Second)
	}
}

func reconnectCluster(options *redis.ClusterOptions) {
	mu.Lock()
	defer mu.Unlock()

	for {
		ctx := context.Background()
		newClient := redis.NewClusterClient(options)
		_, err := newClient.Ping(ctx).Result()
		if err == nil {
			clusterClient = newClient
			log.Info("[REDIS] Successfully reconnected to redis cluster")
			break
		}
		log.Error("[REDIS] Reconnect failed, retrying... " + err.Error())
		time.Sleep(5 * time.Second)
	}
}

func Client() *redis.Client {
	mu.Lock()
	defer mu.Unlock()
	if isCluster {
		return nil
	}
	return client
}

// ClusterClient returns the global Redis cluster client instance
func ClusterClient() *redis.ClusterClient {
	mu.Lock()
	defer mu.Unlock()
	if !isCluster {
		return nil
	}
	return clusterClient
}
