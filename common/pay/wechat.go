package pay

import (
	"context"
	"fmt"
	"time"

	"github.com/go-pay/gopay"
	"github.com/go-pay/gopay/wechat/v3"
	log "github.com/sirupsen/logrus"
)

// ==================== 微信支付配置 ====================

// WechatConfig 微信支付配置
type WechatConfig struct {
	MchID      string // 商户号
	AppID      string // 应用ID (公众号/小程序/APP)
	APIv3Key   string // APIv3 密钥
	SerialNo   string // 商户证书序列号
	PrivateKey string // 商户私钥内容
	NotifyURL  string // 异步通知地址
	IsProd     bool   // 是否生产环境
}

// WechatService 微信支付服务
type WechatService struct {
	client *wechat.ClientV3
	config *WechatConfig
}

// NewWechatService 创建微信支付服务
func NewWechatService(config *WechatConfig) (*WechatService, error) {
	// 创建微信支付 V3 客户端
	client, err := wechat.NewClientV3(config.MchID, config.SerialNo, config.APIv3Key, config.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("create wechat client error: %v", err)
	}

	// 自动验签
	err = client.AutoVerifySign()
	if err != nil {
		return nil, fmt.Errorf("auto verify sign error: %v", err)
	}

	return &WechatService{
		client: client,
		config: config,
	}, nil
}

// ==================== 支付接口 ====================

// WechatOrderRequest 微信支付订单请求
type WechatOrderRequest struct {
	OrderNo     string // 订单号
	Amount      int    // 金额 (分)
	Description string // 商品描述
	ClientIP    string // 客户端 IP
	OpenID      string // 用户 OpenID (JSAPI 必填)
}

// WechatOrderResponse 微信支付订单响应
type WechatOrderResponse struct {
	OrderNo    string `json:"order_no"`
	PrepayID   string `json:"prepay_id"`  // 预支付ID
	CodeURL    string `json:"code_url"`   // 二维码链接 (Native)
	H5URL      string `json:"h5_url"`     // H5 支付链接
	JSAPIData  string `json:"jsapi_data"` // JSAPI 调起参数
	ExpireTime string `json:"expire_time"`
}

// CreateNativePayOrder 创建扫码支付订单 (Native)
func (s *WechatService) CreateNativePayOrder(req *WechatOrderRequest) (*WechatOrderResponse, error) {
	expire := time.Now().Add(30 * time.Minute).Format(time.RFC3339)

	bm := make(gopay.BodyMap)
	bm.Set("appid", s.config.AppID)
	bm.Set("mchid", s.config.MchID)
	bm.Set("description", req.Description)
	bm.Set("out_trade_no", req.OrderNo)
	bm.Set("time_expire", expire)
	bm.Set("notify_url", s.config.NotifyURL)
	bm.SetBodyMap("amount", func(bm gopay.BodyMap) {
		bm.Set("total", req.Amount)
		bm.Set("currency", "CNY")
	})

	resp, err := s.client.V3TransactionNative(context.Background(), bm)
	if err != nil {
		return nil, fmt.Errorf("wechat native paycalback error: %v", err)
	}

	if resp.Code != wechat.Success {
		return nil, fmt.Errorf("wechat error: %s", resp.Error)
	}

	return &WechatOrderResponse{
		OrderNo:    req.OrderNo,
		CodeURL:    resp.Response.CodeUrl,
		ExpireTime: expire,
	}, nil
}

// CreateH5PayOrder 创建 H5 支付订单
func (s *WechatService) CreateH5PayOrder(req *WechatOrderRequest) (*WechatOrderResponse, error) {
	expire := time.Now().Add(30 * time.Minute).Format(time.RFC3339)

	bm := make(gopay.BodyMap)
	bm.Set("appid", s.config.AppID)
	bm.Set("mchid", s.config.MchID)
	bm.Set("description", req.Description)
	bm.Set("out_trade_no", req.OrderNo)
	bm.Set("time_expire", expire)
	bm.Set("notify_url", s.config.NotifyURL)
	bm.SetBodyMap("amount", func(bm gopay.BodyMap) {
		bm.Set("total", req.Amount)
		bm.Set("currency", "CNY")
	})
	bm.SetBodyMap("scene_info", func(bm gopay.BodyMap) {
		bm.Set("payer_client_ip", req.ClientIP)
		bm.SetBodyMap("h5_info", func(bm gopay.BodyMap) {
			bm.Set("type", "Wap")
		})
	})

	resp, err := s.client.V3TransactionH5(context.Background(), bm)
	if err != nil {
		return nil, fmt.Errorf("wechat h5 paycalback error: %v", err)
	}

	if resp.Code != wechat.Success {
		return nil, fmt.Errorf("wechat error: %s", resp.Error)
	}

	return &WechatOrderResponse{
		OrderNo:    req.OrderNo,
		H5URL:      resp.Response.H5Url,
		ExpireTime: expire,
	}, nil
}

// CreateJSAPIPayOrder 创建 JSAPI 支付订单 (公众号/小程序)
func (s *WechatService) CreateJSAPIPayOrder(req *WechatOrderRequest) (*WechatOrderResponse, error) {
	if req.OpenID == "" {
		return nil, fmt.Errorf("openid is required for JSAPI paycalback")
	}

	expire := time.Now().Add(30 * time.Minute).Format(time.RFC3339)

	bm := make(gopay.BodyMap)
	bm.Set("appid", s.config.AppID)
	bm.Set("mchid", s.config.MchID)
	bm.Set("description", req.Description)
	bm.Set("out_trade_no", req.OrderNo)
	bm.Set("time_expire", expire)
	bm.Set("notify_url", s.config.NotifyURL)
	bm.SetBodyMap("amount", func(bm gopay.BodyMap) {
		bm.Set("total", req.Amount)
		bm.Set("currency", "CNY")
	})
	bm.SetBodyMap("payer", func(bm gopay.BodyMap) {
		bm.Set("openid", req.OpenID)
	})

	resp, err := s.client.V3TransactionJsapi(context.Background(), bm)
	if err != nil {
		return nil, fmt.Errorf("wechat jsapi paycalback error: %v", err)
	}

	if resp.Code != wechat.Success {
		return nil, fmt.Errorf("wechat error: %s", resp.Error)
	}

	// 生成 JSAPI 调起参数
	jsapi, err := s.client.PaySignOfJSAPI(s.config.AppID, resp.Response.PrepayId)
	if err != nil {
		return nil, fmt.Errorf("generate jsapi sign error: %v", err)
	}

	return &WechatOrderResponse{
		OrderNo:    req.OrderNo,
		PrepayID:   resp.Response.PrepayId,
		JSAPIData:  jsapi.PaySign, // 只返回签名字符串
		ExpireTime: expire,
	}, nil
}

// ==================== 订单查询 ====================

// WechatQueryResult 查询结果
type WechatQueryResult struct {
	OrderNo        string `json:"order_no"`
	TransactionID  string `json:"transaction_id"` // 微信支付订单号
	TradeState     string `json:"trade_state"`
	TradeStateDesc string `json:"trade_state_desc"`
	TotalAmount    int    `json:"total_amount"` // 分
	PayerTotal     int    `json:"payer_total"`  // 用户实付金额
	IsPaid         bool   `json:"is_paid"`
}

// QueryOrder 查询订单状态
func (s *WechatService) QueryOrder(orderNo string) (*WechatQueryResult, error) {
	resp, err := s.client.V3TransactionQueryOrder(context.Background(), wechat.OutTradeNo, orderNo)
	if err != nil {
		return nil, fmt.Errorf("wechat query order error: %v", err)
	}

	if resp.Code != wechat.Success {
		return nil, fmt.Errorf("wechat error: %s", resp.Error)
	}

	isPaid := resp.Response.TradeState == "SUCCESS"

	result := &WechatQueryResult{
		OrderNo:        resp.Response.OutTradeNo,
		TransactionID:  resp.Response.TransactionId,
		TradeState:     resp.Response.TradeState,
		TradeStateDesc: resp.Response.TradeStateDesc,
		IsPaid:         isPaid,
	}

	if resp.Response.Amount != nil {
		result.TotalAmount = resp.Response.Amount.Total
		result.PayerTotal = resp.Response.Amount.PayerTotal
	}

	return result, nil
}

// ==================== 回调验签 ====================

// WechatNotifyResult 通知结果
type WechatNotifyResult struct {
	OrderNo       string
	TransactionID string
	TradeState    string
	TotalAmount   int
	PayerTotal    int
	SuccessTime   string
}

// VerifyNotify 验证异步通知
func (s *WechatService) VerifyNotify(notifyReq *wechat.V3NotifyReq) (*WechatNotifyResult, error) {
	// 解密通知内容
	result, err := notifyReq.DecryptPayCipherText(s.config.APIv3Key)
	if err != nil {
		return nil, fmt.Errorf("decrypt notify error: %v", err)
	}

	notify := &WechatNotifyResult{
		OrderNo:       result.OutTradeNo,
		TransactionID: result.TransactionId,
		TradeState:    result.TradeState,
		SuccessTime:   result.SuccessTime,
	}

	if result.Amount != nil {
		notify.TotalAmount = result.Amount.Total
		notify.PayerTotal = result.Amount.PayerTotal
	}

	return notify, nil
}

// ==================== 关闭订单 ====================

// CloseOrder 关闭订单
func (s *WechatService) CloseOrder(orderNo string) error {
	resp, err := s.client.V3TransactionCloseOrder(context.Background(), orderNo)
	if err != nil {
		return fmt.Errorf("wechat close order error: %v", err)
	}

	if resp.Code != wechat.Success {
		return fmt.Errorf("wechat error: %s", resp.Error)
	}

	return nil
}

// ==================== 全局实例 ====================

var wechatService *WechatService

// InitWechat 初始化微信支付服务
func InitWechat(config *WechatConfig) error {
	var err error
	wechatService, err = NewWechatService(config)
	if err != nil {
		return err
	}
	log.Info("[Wechat] Payment service initialized")
	return nil
}

// GetWechatService 获取微信支付服务实例
func GetWechatService() *WechatService {
	return wechatService
}
