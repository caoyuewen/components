package dbmysql

import (
	"fmt"
	"sync/atomic"
	"time"

	log "github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// dbc 数据库连接单例（使用 atomic.Pointer 实现无锁并发安全）
var (
	dbc   atomic.Pointer[gorm.DB]
	info  MysqlInfo
	check time.Duration
)

type MysqlInfo struct {
	Address  string
	User     string
	Password string
	DBName   string
	Logger   logger.Interface
}

// StartUp 在中间件中初始化 MySQL 连接
func StartUp(msqlInfo MysqlInfo, checkInterval time.Duration) {
	check = checkInterval
	info = msqlInfo
	connDB()
	go checkConnection()
}

func connDB() {
	// 数据源名称
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		info.User,
		info.Password,
		info.Address,
		info.DBName,
	)

	// 使用新版的 gorm 连接 MySQL
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: info.Logger})
	if err != nil {
		log.Errorf("[MYSQL] MySQL database connection failed: %v", err)
		return
	}
	log.Infof("[MYSQL] MySQL connected successfully: %s", info.Address)

	// 获取通用数据库对象 sql.DB 以便设置连接池等
	sqlDB, err := db.DB()
	if err != nil {
		log.Errorf("[MYSQL] Failed to get generic database object: %v", err)
		return
	}

	// 设置连接池
	sqlDB.SetMaxIdleConns(20)                  // 空闲连接数
	sqlDB.SetMaxOpenConns(100)                 // 最大打开连接数
	sqlDB.SetConnMaxLifetime(30 * time.Second) // 连接的最大生命周期

	// 原子存储，无锁
	dbc.Store(db)
}

// 定期检查数据库连接
func checkConnection() {
	ticker := time.NewTicker(check)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			if !checkDBConnection() {
				reconnectDB()
			}
		}
	}
}

// 检测数据库连接是否断开
func checkDBConnection() bool {
	db := dbc.Load()
	if db == nil {
		return false
	}
	sqlDB, err := db.DB()
	if err != nil {
		log.Errorf("[MYSQL] Failed to get generic database object: %v", err)
		return false
	}
	err = sqlDB.Ping()
	if err != nil {
		log.Errorf("[MYSQL] MySQL database connection lost: %v", err)
		return false
	}
	return true
}

// 重连数据库
func reconnectDB() {
	log.Info("[MYSQL] Attempting to reconnect to the MySQL database...")
	connDB()
	if dbc.Load() != nil {
		log.Info("[MYSQL] Successfully reconnected to the MySQL database")
	}
}

// Client 获取数据库客户端实例（原子读取，无锁）
func Client() *gorm.DB {
	return dbc.Load()
}
