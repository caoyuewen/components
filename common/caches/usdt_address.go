package caches

import (
	"context"
	"errors"

	"github.com/caoyuewen/components/dbs/dbredis"
	"github.com/redis/go-redis/v9"
	"github.com/shopspring/decimal"
	log "github.com/sirupsen/logrus"
)

// Redis Key
const (
	RedisKeyUsdtAddressPool  = "u_pool" // 地址池 (List)
	RedisKeyUsdtAddressOrder = "u:"     // 地址订单占位 (ZSet) + address
)

var UsdtAddress usdtAddress

type usdtAddress struct{}

type UsdtAddressOrderInfo struct {
	OrderId string
	Amount  string
}

// FlushPool 刷新地址池到 Redis
func (u *usdtAddress) FlushPool(addrList []string) error {
	ctx := context.Background()
	rdb := dbredis.Client()

	// 清空 Redis 缓存并写入新数据 (List)
	pipe := rdb.TxPipeline()
	pipe.Del(ctx, RedisKeyUsdtAddressPool)

	if len(addrList) == 0 {
		log.Warn("Warning: UsdtAddress pool is empty")
		_, err := pipe.Exec(ctx)
		if err != nil {
			log.Info("UsdtAddressFlushAll err:", err)
			return err
		}
		return nil
	}

	// 转换为 interface{} 切片
	args := make([]interface{}, len(addrList))
	for i, addr := range addrList {
		args[i] = addr
	}
	pipe.RPush(ctx, RedisKeyUsdtAddressPool, args...)

	_, err := pipe.Exec(ctx)
	if err != nil {
		log.Error("UsdtAddressFlushAll err:", err)
		return err
	}
	log.Info("UsdtAddressFlushAll success, count:", len(addrList))

	return nil
}

// PoolCount 获取缓存中的地址总条数
func (u *usdtAddress) PoolCount() (int64, error) {
	ctx := context.Background()
	return dbredis.Client().LLen(ctx, RedisKeyUsdtAddressPool).Result()
}

// Pop 从地址池中分配一个可用地址 (检查金额冲突)
func (u *usdtAddress) Pop(amount decimal.Decimal) (string, error) {

	ctx := context.Background()
	rdb := dbredis.Client()

	count, err := u.PoolCount()
	if err != nil {
		return "", err
	}

	for i := 0; i < int(count); i++ {
		// 从地址池循环取一个地址 (RPopLPush 循环)
		addr, err := rdb.RPopLPush(ctx, RedisKeyUsdtAddressPool, RedisKeyUsdtAddressPool).Result()
		if errors.Is(err, redis.Nil) {
			return "", errors.New("address pool is empty")
		} else if err != nil {
			return "", err
		}

		// 针对当前地址，构造订单 ZSet key
		key := RedisKeyUsdtAddressOrder + addr
		minAmount := amount.Sub(decimal.NewFromInt(2))
		maxAmount := amount.Add(decimal.NewFromInt(2))

		// ZSet：检查是否存在区间内的金额 (防止金额冲突)
		cnt, err := rdb.ZCount(ctx, key, minAmount.String(), maxAmount.String()).Result()
		if err != nil {
			return "", err
		}

		if cnt > 0 {
			// 有冲突，换下一个地址
			orders, _ := u.GetOrders(addr, amount)
			log.Infof("%s Pop aleadly use check next: %+v", addr, orders)
			continue
		}

		return addr, nil
	}

	return "", errors.New("no available address")
}

// SetOrder 设置地址订单占位
func (u *usdtAddress) SetOrder(orderId, addr string, amount decimal.Decimal) error {
	ctx := context.Background()
	key := RedisKeyUsdtAddressOrder + addr
	fScore, _ := amount.Float64()
	rdb := dbredis.Client()

	return rdb.ZAdd(ctx, key, redis.Z{
		Score:  fScore,
		Member: orderId,
	}).Err()
}

// GetOrders 根据地址和金额范围查询订单占位
func (u *usdtAddress) GetOrders(addr string, amount decimal.Decimal) ([]UsdtAddressOrderInfo, error) {
	ctx := context.Background()
	key := RedisKeyUsdtAddressOrder + addr

	opt := &redis.ZRangeBy{
		Min: amount.Sub(decimal.NewFromInt(2)).String(),
		Max: amount.Add(decimal.NewFromInt(2)).String(),
	}

	orders, err := dbredis.Client().ZRangeByScoreWithScores(ctx, key, opt).Result()
	if err != nil {
		return nil, err
	}

	var res []UsdtAddressOrderInfo
	for _, o := range orders {
		a := decimal.NewFromFloat(o.Score).Round(2)
		res = append(res, UsdtAddressOrderInfo{
			OrderId: o.Member.(string),
			Amount:  a.String(),
		})
	}

	return res, nil
}

// DelOrder 删除地址订单占位
func (u *usdtAddress) DelOrder(addr, orderId string) error {
	ctx := context.Background()
	key := RedisKeyUsdtAddressOrder + addr

	_, err := dbredis.Client().ZRem(ctx, key, orderId).Result()
	if err != nil {
		log.Error("DelOrder err:", err.Error())
		return err
	}
	log.Info("DelOrder success", addr, orderId)
	return nil
}
