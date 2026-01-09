package models

import (
	"fmt"

	"github.com/caoyuewen/components/dbs/dbmysql"
	"github.com/fbsobreira/gotron-sdk/pkg/address"
	"gorm.io/gorm/clause"
)

// ValidateTRC20Address 验证 TRC20 地址 (使用 gotron-sdk)
func ValidateTRC20Address(addr string) error {

	_, err := address.Base58ToAddress(addr)
	if err != nil {
		return fmt.Errorf("无效的 TRC20 地址: %v", err)
	}

	return nil
}

var UsdtAddressRepo = dbmysql.NewBaseRepository[UsdtAddress]("id")

// UsdtAddress 充值地址
type UsdtAddress struct {
	ID            string `json:"id" gorm:"primaryKey;type:varchar(64)"`
	Address       string `json:"address" gorm:"type:varchar(100);uniqueIndex;not null"`
	IsActive      int    `json:"is_active" gorm:"column:is_active;type:tinyint;default:1;not null"` // 1 启用 2 停用
	Priority      int    `json:"priority"  gorm:"type:int;default:1;not null"`
	TotalReceived string `json:"total_received" gorm:"-"` // 累计收款 (计算字段，不存数据库)
	CreatedAt     int64  `json:"created_at" gorm:"autoCreateTime:milli"`
	UpdatedAt     int64  `json:"updated_at" gorm:"autoUpdateTime:milli"`
}

func (*UsdtAddress) TableName() string { return "usdt_address" }

// UsdtAddressTotalReceived 统计地址累计收款金额
func UsdtAddressTotalReceived(address string) string {

	var total struct {
		Sum float64
	}
	// 从订单表统计已支付的金额
	dbmysql.Client().Model(&GoodsOrder{}).
		Select("COALESCE(SUM(CAST(real_amount AS DECIMAL(20,8))), 0) as sum").
		Where("to_address = ? AND order_status = 2", address). // 2 = 已支付
		Scan(&total)

	if total.Sum == 0 {
		return "0"
	}
	return fmt.Sprintf("%.2f", total.Sum)
}

// UsdtAddressActiveList 可用的地址
func UsdtAddressActiveList() ([]string, error) {

	var res []string

	// 查询数据库全部启用的地址
	cond := []interface{}{
		clause.Eq{Column: "is_active", Value: 1},
	}

	find, err := UsdtAddressRepo.Find("priority desc", cond...)
	if err != nil {
		return nil, err
	}

	for _, v := range find {
		res = append(res, v.Address)
	}

	return res, nil

}
