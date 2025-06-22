package gen

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/bwmarrin/snowflake"
)

var (
	once    sync.Once
	node    *snowflake.Node
	nodeID  int64 = 1
	initErr error
)

// 设置 Snowflake 起始时间
func initNode() {
	snowflake.Epoch = time.Date(1995, 9, 5, 2, 2, 95, 95, time.UTC).UnixNano() / 1e6
	node, initErr = snowflake.NewNode(nodeID)
	if initErr != nil {
		log.Printf("Snowflake node init failed: %v", initErr)
	}
}

// Init 可供外部主动初始化使用
func Init(id int64) {
	nodeID = id
	once.Do(initNode)
}

func IdInt64() int64 {
	once.Do(initNode)
	if initErr != nil || node == nil {
		return 0
	}
	return node.Generate().Int64()
}

func IdString() string {
	return fmt.Sprintf("%d", IdInt64())
}
