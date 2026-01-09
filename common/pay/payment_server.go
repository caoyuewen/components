package pay

import (
	"errors"
	"fmt"
)

const (
	OrderStatusPending = 1 // 待支付
	OrderStatusSuccess = 2 // 成功
	OrderStatusFailed  = 3 // 失败
	OrderStatusExpired = 4 // 已过期
)

const (
	PayTypeUsdt   = 1 // USDT
	PayTypeAlipay = 2 // Alipay
	PayTypeWechat = 3 // WeiChat
)

var PayTypeMap = map[int]string{
	PayTypeUsdt:   "USDT",
	PayTypeAlipay: "Alipay",
	PayTypeWechat: "WeiChat",
}

var PayTypeSymbolMap = map[int]string{
	PayTypeUsdt:   "$",
	PayTypeAlipay: "¥",
	PayTypeWechat: "¥",
}

var OrderStatusMap = map[int]string{
	OrderStatusPending: "待支付",
	OrderStatusSuccess: "成功",
	OrderStatusFailed:  "失败",
	OrderStatusExpired: "已过期",
}

// PaymentService 所有的第三方支付渠道必须实现以下接口
type PaymentService interface {
	CallDeposit(id, amount string) (CallDepositResult, error)                               // 调用三方充值请求
	CallDepositOrderQuery(orderID, externalOrderID string) (PaymentOrderQueryResult, error) // 查询三方充值订单
}

type Payment struct {
	Name        string
	PayService  PaymentService
	PaymentType int
}

// PaymentMap 支付渠道映射 key = 三方渠道名 ; v = 对应第三方渠道
// 只允许启动时加载 运行时只读
var PaymentMap = map[string]Payment{}

func paymentRegister(payment Payment) {
	PaymentMap[payment.Name] = payment
}

type PaymentServerCondition struct {
	Name   string
	Amount string
	Config string
}

func CreatePaymentServer(c PaymentServerCondition) (PaymentService, error) {

	var (
		paymentService PaymentService
	)

	ps, ok := PaymentMap[c.Name]
	if !ok {
		errStr := fmt.Sprintf("not found payment factory by name:%s", c.Name)
		return paymentService, errors.New(errStr)
	}

	return ps.PayService, nil
}

// CallDepositResult Deposit调用第三方成功的通用返回
type CallDepositResult struct {
	ExternalOrderID string // 第三方返回的id
	ToAddress       string // 收款的区块链地址
	PayUrl          string // 支付url
}

// PaymentOrderQueryResult 查询三方订单的通用返回
type PaymentOrderQueryResult struct {
	OrderNo         string `json:"order_no"`          // 我方订单号
	ExternalOrderID string `json:"external_order_id"` // 三方订单号
	TxId            string `json:"tx_id"`             // 链上交易hash
	ToAddress       string `json:"to_address"`        // 收款地址
	FromAddress     string `json:"from_address"`      // 付款地址
	Status          int    `json:"status"`            // 我们维护的订单状态
	ExternalStatus  string `json:"external_status"`   // 三方的订单状态(三方返回的什么 就存什么)
	PayAt           int64  `json:"pay_at"`            // 真实付款时间
	Amount          string `json:"amount"`            // 订单金额
	RealAmount      string `json:"real_amount"`       // 真实收到的金额
	Fee             string `json:"fee"`               // 税收
}
