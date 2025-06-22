package dbmongo

import (
	"context"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client *mongo.Client
	once   sync.Once
	mu     sync.Mutex
)

type MongoInfo struct {
	Address  string
	User     string
	Password string
}

// StartUp initializes the global MongoDB client or cluster client with the provided options and starts the connection checker
func StartUp(info MongoInfo, checkInterval time.Duration) {
	uri := "mongodb://" + info.User + ":" + info.Password + "@" + info.Address
	StartUpByUri(uri, checkInterval)
}

// StartUpByUri initializes the global MongoDB client or cluster client with the provided options and starts the connection checker
func StartUpByUri(uri string, checkInterval time.Duration) {
	var err error
	once.Do(func() {
		client, err = mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
		if err != nil {
			log.Error("[MONGO] Failed to initialize MongoDB: " + err.Error())
			return
		}
		if err = client.Ping(context.Background(), nil); err != nil {
			log.Error("[MONGO] Failed to ping MongoDB: " + err.Error())
			return
		}
		log.Info("[MONGO] MongoDB initialized success")
		go connectionChecker(uri, checkInterval)

	})
	return
}

// connectionChecker periodically checks the MongoDB connection and reconnects if necessary
func connectionChecker(uri string, checkInterval time.Duration) {
	ticker := time.NewTicker(checkInterval)
	defer ticker.Stop()

	for range ticker.C {
		if err := client.Ping(context.Background(), nil); err != nil {
			log.Warn("[MONGO] MongoDB connection lost, attempting to reconnect...")
			reconnect(uri)
		}
	}
}

// reconnect attempts to reconnect to MongoDB
func reconnect(uri string) {
	mu.Lock()
	defer mu.Unlock()

	for {
		newClient, err := mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
		if err == nil {
			if err = newClient.Ping(context.Background(), nil); err == nil {
				client = newClient
				log.Info("[MONGO] Successfully reconnected to MongoDB")
				break
			}
		}
		log.Error("[MONGO] Reconnect failed, retrying... " + err.Error())
		time.Sleep(5 * time.Second) // wait before retrying
	}
}

// Client returns the global MongoDB client instance
func Client() *mongo.Client {
	mu.Lock()
	defer mu.Unlock()
	return client
}
