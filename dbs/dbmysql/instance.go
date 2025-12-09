package dbmysql

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"time"
)

// dbc 数据库链接单例
var (
	dbc  *gorm.DB
	info MysqlInfo
)

var check time.Duration

type MysqlInfo struct {
	Address  string
	User     string
	Password string
	DBName   string
	Logger   logger.Interface
}

// StartUp 在中间件中初始化mysql链接
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

	// 使用新版的gorm连接MySQL
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: info.Logger})
	//db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Errorf("[MYSQL] Mysql Database connection failed: %v", err)
		return
	}
	log.Infof("[MYSQL] Mysql connected successfully: %s", info.Address)

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
	dbc = db
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
	sqlDB, err := dbc.DB()
	if err != nil {
		log.Errorf("[MYSQL] Failed to get generic database object: %v", err)
		return false
	}
	err = sqlDB.Ping()
	if err != nil {
		log.Errorf("[MYSQL] Mysql Database connection lost: %v", err)
		return false
	}
	return true
}

// 重连数据库
func reconnectDB() {
	log.Info("[MYSQL] Attempting to reconnect to the MySQL database...")
	connDB()
	if dbc != nil {
		log.Info("[MYSQL] Successfully reconnected to the MySQL database")
	}
}

// 获取数据库客户端实例
func Client() *gorm.DB {
	return dbc
}
