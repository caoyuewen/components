package dbmysql

import (
	"context"
	"fmt"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// 数据库连接单例
var (
	dbc   *gorm.DB
	info  MysqlInfo
	check time.Duration
	once  sync.Once
	mu    sync.RWMutex
)

// MysqlInfo MySQL 配置信息
type MysqlInfo struct {
	Address  string          // 数据库地址 host:port
	User     string          // 用户名
	Password string          // 密码
	DBName   string          // 数据库名
	Logger   logger.Interface // 日志接口

	// 连接池配置（可选）
	MaxIdleConns    int           // 最大空闲连接数，默认 20
	MaxOpenConns    int           // 最大打开连接数，默认 100
	ConnMaxLifetime time.Duration // 连接最大生命周期，默认 30s
	ConnMaxIdleTime time.Duration // 空闲连接最大生命周期，默认 10m
}

// StartUp 初始化 MySQL 连接
func StartUp(msqlInfo MysqlInfo, checkInterval time.Duration) {
	once.Do(func() {
		check = checkInterval
		info = msqlInfo
		setDefaultConfig()
		connDB()
		go checkConnection()
	})
}

// setDefaultConfig 设置默认配置
func setDefaultConfig() {
	if info.MaxIdleConns <= 0 {
		info.MaxIdleConns = 20
	}
	if info.MaxOpenConns <= 0 {
		info.MaxOpenConns = 100
	}
	if info.ConnMaxLifetime <= 0 {
		info.ConnMaxLifetime = 30 * time.Second
	}
	if info.ConnMaxIdleTime <= 0 {
		info.ConnMaxIdleTime = 10 * time.Minute
	}
}

// connDB 连接数据库
func connDB() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local&timeout=10s&readTimeout=30s&writeTimeout=60s",
		info.User,
		info.Password,
		info.Address,
		info.DBName,
	)

	config := &gorm.Config{
		Logger:                                   info.Logger,
		SkipDefaultTransaction:                   true,  // 跳过默认事务，提升性能
		PrepareStmt:                              true,  // 缓存预编译语句
		DisableForeignKeyConstraintWhenMigrating: true,  // 迁移时禁用外键约束
	}

	db, err := gorm.Open(mysql.Open(dsn), config)
	if err != nil {
		log.Errorf("[MYSQL] Database connection failed: %v", err)
		return
	}
	log.Infof("[MYSQL] Connected successfully: %s/%s", info.Address, info.DBName)

	// 获取底层 sql.DB 设置连接池
	sqlDB, err := db.DB()
	if err != nil {
		log.Errorf("[MYSQL] Failed to get sql.DB: %v", err)
		return
	}

	// 设置连接池参数
	sqlDB.SetMaxIdleConns(info.MaxIdleConns)
	sqlDB.SetMaxOpenConns(info.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(info.ConnMaxLifetime)
	sqlDB.SetConnMaxIdleTime(info.ConnMaxIdleTime)

	mu.Lock()
	dbc = db
	mu.Unlock()
}

// checkConnection 定期检查数据库连接
func checkConnection() {
	ticker := time.NewTicker(check)
	defer ticker.Stop()
	for range ticker.C {
		if !isConnected() {
			reconnectDB()
		}
	}
}

// isConnected 检测数据库连接是否正常
func isConnected() bool {
	mu.RLock()
	db := dbc
	mu.RUnlock()

	if db == nil {
		return false
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Errorf("[MYSQL] Failed to get sql.DB: %v", err)
		return false
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := sqlDB.PingContext(ctx); err != nil {
		log.Errorf("[MYSQL] Connection lost: %v", err)
		return false
	}
	return true
}

// reconnectDB 重连数据库
func reconnectDB() {
	log.Info("[MYSQL] Attempting to reconnect...")
	connDB()

	mu.RLock()
	connected := dbc != nil
	mu.RUnlock()

	if connected {
		log.Info("[MYSQL] Reconnected successfully")
	}
}

// Client 获取数据库客户端实例
func Client() *gorm.DB {
	mu.RLock()
	defer mu.RUnlock()
	return dbc
}

// ClientWithContext 获取带上下文的数据库客户端
func ClientWithContext(ctx context.Context) *gorm.DB {
	mu.RLock()
	defer mu.RUnlock()
	if dbc == nil {
		return nil
	}
	return dbc.WithContext(ctx)
}

// IsConnected 检查连接是否正常
func IsConnected() bool {
	return isConnected()
}

// Close 关闭数据库连接
func Close() error {
	mu.Lock()
	defer mu.Unlock()

	if dbc == nil {
		return nil
	}

	sqlDB, err := dbc.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// Stats 获取连接池统计信息
func Stats() map[string]interface{} {
	mu.RLock()
	db := dbc
	mu.RUnlock()

	if db == nil {
		return nil
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil
	}

	stats := sqlDB.Stats()
	return map[string]interface{}{
		"max_open_connections": stats.MaxOpenConnections,
		"open_connections":     stats.OpenConnections,
		"in_use":               stats.InUse,
		"idle":                 stats.Idle,
		"wait_count":           stats.WaitCount,
		"wait_duration":        stats.WaitDuration.String(),
		"max_idle_closed":      stats.MaxIdleClosed,
		"max_lifetime_closed":  stats.MaxLifetimeClosed,
	}
}
