package models

import (
	"github.com/caoyuewen/components/dbs/dbmysql"
)

var GoodsOrderRepo = dbmysql.NewBaseRepository[GoodsOrder]("id")

const (
	OrderExpiredTime = 15 // 单位分钟
)

type GoodsOrder struct {
	ID              string `gorm:"primaryKey;not null" json:"id"`                    // 订单 ID / 订单号
	GoodsId         string `gorm:"not null;index" json:"goods_id"`                   // 商品 ID / VIP 套餐 ID
	GoodsTitle      string `gorm:"type:VARCHAR(255);not null" json:"goods_title"`    // 商品标题
	Uid             string `gorm:"not null;index" json:"uid"`                        // 用户 ID
	UserEmail       string `gorm:"type:varchar(128);" json:"user_email"`             // 用户邮箱（冗余字段，方便查询）
	ChannelId       string `gorm:"type:BIGINT UNSIGNED" json:"channel_id" `          // 三方渠道唯一标识
	ChannelName     string `gorm:"type:varchar(50);" json:"channel_name"`            // 三方渠道名称
	PayType         int    `gorm:"type:int;default:1" json:"pay_type"`               // 支付方式 1 USDT 2 支付宝 3 微信
	Amount          string `gorm:"type:decimal(20,8);not null" json:"amount"`        // 应付金额
	RealAmount      string `gorm:"type:decimal(20,8);not null" json:"real_amount"`   // 实际到账金额
	ExternalOrderId string `gorm:"type:varchar(100);index" json:"external_order_id"` // 三方渠道订单号
	OrderStatus     int    `gorm:"type:tinyint;not null;index" json:"order_status"`  // 订单状态
	ExternalStatus  string `gorm:"type:varchar(50);" json:"external_status"`         // 三方返回的订单状态

	// USDT 支付相关字段
	FromAddress string `gorm:"type:varchar(100)" json:"from_address"`   // 付款方地址
	ToAddress   string `gorm:"type:varchar(100)" json:"to_address"`     // 收款方地址
	TxHash      string `gorm:"type:varchar(100)" json:"tx_hash"`        // 交易哈希
	ExpireTime  int64  `gorm:"type:BIGINT;not null" json:"expire_time"` // 过期时间

	PaidAt     int64  `gorm:"type:BIGINT;not null" json:"paid_at"`   // 实际付款时间
	FailReason string `gorm:"type:varchar(255);" json:"fail_reason"` // 失败原因
	Remark     string `gorm:"type:varchar(255);" json:"remark"`      // 备注
	CreatedAt  int64  `gorm:"type:BIGINT;not null" json:"created_at"`
	UpdatedAt  int64  `gorm:"type:BIGINT;not null" json:"updated_at"`

	// 返回给前端的字段非数据库字段
	OrderStatusName string `gorm:"-" json:"order_status_name"` // 订单状态
	PayTypeName     string `gorm:"-" json:"pay_type_name"`     // 支付方式
	AmountSymbol    string `gorm:"-" json:"amount_symbol"`     // 金额单位符号

}

func (*GoodsOrder) TableName() string { return "goods_order" }
