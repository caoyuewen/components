package pay

import (
	"context"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-pay/gopay"
	"github.com/go-pay/gopay/alipay"
	log "github.com/sirupsen/logrus"
)

// ==================== 支付宝配置 ====================

// AlipayConfig 支付宝配置
type AlipayConfig struct {
	AppID           string // 应用ID
	PrivateKey      string // 应用私钥 (RSA2)
	AlipayPublicKey string // 支付宝公钥 (用于验签)
	NotifyURL       string // 异步通知地址
	ReturnURL       string // 同步跳转地址
	IsProd          bool   // 是否生产环境
}

// AlipayService 支付宝支付服务
type AlipayService struct {
	client *alipay.Client
	config *AlipayConfig
}

// NewAlipayService 创建支付宝服务
func NewAlipayService(config *AlipayConfig) (*AlipayService, error) {
	// 创建支付宝客户端
	client, err := alipay.NewClient(config.AppID, config.PrivateKey, config.IsProd)
	if err != nil {
		return nil, fmt.Errorf("create alipay client error: %v", err)
	}

	// 设置支付宝公钥证书 (用于验签)
	client.SetAliPayPublicCertSN(config.AlipayPublicKey)

	return &AlipayService{
		client: client,
		config: config,
	}, nil
}

// ==================== 支付接口 ====================

// AlipayOrderRequest 支付宝订单请求
type AlipayOrderRequest struct {
	OrderNo     string  // 订单号
	Amount      float64 // 金额 (元)
	Subject     string  // 商品标题
	Description string  // 商品描述
}

// AlipayOrderResponse 支付宝订单响应
type AlipayOrderResponse struct {
	OrderNo    string `json:"order_no"`
	PayURL     string `json:"pay_url"`     // 支付链接 (PC 网页支付)
	QRCodeURL  string `json:"qr_code_url"` // 二维码内容 (扫码支付)
	H5URL      string `json:"h5_url"`      // H5 支付链接
	ExpireTime string `json:"expire_time"`
}

// CreatePCPayOrder 创建 PC 网页支付订单
func (s *AlipayService) CreatePCPayOrder(req *AlipayOrderRequest) (*AlipayOrderResponse, error) {
	bm := make(gopay.BodyMap)
	bm.Set("subject", req.Subject)
	bm.Set("out_trade_no", req.OrderNo)
	bm.Set("total_amount", fmt.Sprintf("%.2f", req.Amount))
	bm.Set("product_code", "FAST_INSTANT_TRADE_PAY")

	// 设置超时时间 30 分钟
	bm.Set("time_expire", time.Now().Add(30*time.Minute).Format("2006-01-02 15:04:05"))

	// 发起请求
	payUrl, err := s.client.TradePagePay(context.Background(), bm)
	if err != nil {
		return nil, fmt.Errorf("alipay trade page paycalback error: %v", err)
	}

	return &AlipayOrderResponse{
		OrderNo: req.OrderNo,
		PayURL:  payUrl,
	}, nil
}

// CreateQRCodePayOrder 创建扫码支付订单
func (s *AlipayService) CreateQRCodePayOrder(req *AlipayOrderRequest) (*AlipayOrderResponse, error) {
	bm := make(gopay.BodyMap)
	bm.Set("subject", req.Subject)
	bm.Set("out_trade_no", req.OrderNo)
	bm.Set("total_amount", fmt.Sprintf("%.2f", req.Amount))

	// 发起预创建请求
	resp, err := s.client.TradePrecreate(context.Background(), bm)
	if err != nil {
		return nil, fmt.Errorf("alipay trade precreate error: %v", err)
	}

	if resp.Response.Code != "10000" {
		return nil, fmt.Errorf("alipay error: %s - %s", resp.Response.Code, resp.Response.Msg)
	}

	return &AlipayOrderResponse{
		OrderNo:   req.OrderNo,
		QRCodeURL: resp.Response.QrCode, // 二维码内容
	}, nil
}

// CreateH5PayOrder 创建 H5 支付订单 (手机浏览器)
func (s *AlipayService) CreateH5PayOrder(req *AlipayOrderRequest) (*AlipayOrderResponse, error) {
	bm := make(gopay.BodyMap)
	bm.Set("subject", req.Subject)
	bm.Set("out_trade_no", req.OrderNo)
	bm.Set("total_amount", fmt.Sprintf("%.2f", req.Amount))
	bm.Set("product_code", "QUICK_WAP_WAY")
	bm.Set("quit_url", s.config.ReturnURL)

	// 发起请求
	payUrl, err := s.client.TradeWapPay(context.Background(), bm)
	if err != nil {
		return nil, fmt.Errorf("alipay trade wap paycalback error: %v", err)
	}

	return &AlipayOrderResponse{
		OrderNo: req.OrderNo,
		H5URL:   payUrl,
	}, nil
}

// ==================== 订单查询 ====================

// AlipayQueryResult 查询结果
type AlipayQueryResult struct {
	OrderNo        string `json:"order_no"`
	TradeNo        string `json:"trade_no"` // 支付宝交易号
	TradeStatus    string `json:"trade_status"`
	TotalAmount    string `json:"total_amount"`
	BuyerPayAmount string `json:"buyer_pay_amount"`
	IsPaid         bool   `json:"is_paid"`
}

// QueryOrder 查询订单状态
func (s *AlipayService) QueryOrder(orderNo string) (*AlipayQueryResult, error) {
	bm := make(gopay.BodyMap)
	bm.Set("out_trade_no", orderNo)

	resp, err := s.client.TradeQuery(context.Background(), bm)
	if err != nil {
		return nil, fmt.Errorf("alipay trade query error: %v", err)
	}

	if resp.Response.Code != "10000" {
		return nil, fmt.Errorf("alipay error: %s - %s", resp.Response.Code, resp.Response.Msg)
	}

	isPaid := resp.Response.TradeStatus == "TRADE_SUCCESS" || resp.Response.TradeStatus == "TRADE_FINISHED"

	return &AlipayQueryResult{
		OrderNo:        resp.Response.OutTradeNo,
		TradeNo:        resp.Response.TradeNo,
		TradeStatus:    resp.Response.TradeStatus,
		TotalAmount:    resp.Response.TotalAmount,
		BuyerPayAmount: resp.Response.BuyerPayAmount,
		IsPaid:         isPaid,
	}, nil
}

// ==================== 回调验签 ====================

// VerifyNotify 验证异步通知签名
func (s *AlipayService) VerifyNotify(c *gin.Context) (gopay.BodyMap, error) {
	// 解析通知内容
	notifyReq, err := alipay.ParseNotifyToBodyMap(c.Request)
	if err != nil {
		return nil, fmt.Errorf("parse notify error: %v", err)
	}

	// 验证签名 (使用公钥证书)
	ok, err := alipay.VerifySignWithCert(s.config.AlipayPublicKey, notifyReq)
	if err != nil {
		return nil, fmt.Errorf("verify sign error: %v", err)
	}
	if !ok {
		return nil, fmt.Errorf("sign verify failed")
	}

	return notifyReq, nil
}

// ==================== 关闭订单 ====================

// CloseOrder 关闭订单
func (s *AlipayService) CloseOrder(orderNo string) error {
	bm := make(gopay.BodyMap)
	bm.Set("out_trade_no", orderNo)

	resp, err := s.client.TradeClose(context.Background(), bm)
	if err != nil {
		return fmt.Errorf("alipay trade close error: %v", err)
	}

	if resp.Response.Code != "10000" {
		return fmt.Errorf("alipay error: %s - %s", resp.Response.Code, resp.Response.Msg)
	}

	return nil
}

// ==================== 全局实例 ====================

var alipayService *AlipayService

// InitAlipay 初始化支付宝服务
func InitAlipay(config *AlipayConfig) error {
	var err error
	alipayService, err = NewAlipayService(config)
	if err != nil {
		return err
	}
	log.Info("[Alipay] Service initialized")
	return nil
}

// GetAlipayService 获取支付宝服务实例
func GetAlipayService() *AlipayService {
	return alipayService
}
