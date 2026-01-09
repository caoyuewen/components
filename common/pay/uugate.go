package pay

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/shopspring/decimal"
)

var (
	uugateDevCfg  = ""
	uugateProdCfg = ""
)

//type UugateFactory struct{}
//
//func (f *UugateFactory) Create(config string) (PaymentService, error) {
//	var uugate Uugate
//	if config == "" || config == "{}" {
//		uugate.Uid = "893675"
//		uugate.ApiKey = "f9d74b0d05863aa20f70cd0ea1294f4f"
//		uugate.Domain = "https://open.uugate.com"
//		uugate.JumpUrl = app.Domain() + "/page/user/profile"
//		uugate.EffectiveDuration = 60 * 30
//		return &uugate, nil
//	}
//
//	err := json.Unmarshal([]byte(config), &uugate)
//	if err != nil {
//		return nil, err
//	}
//
//	return &uugate, nil
//}

var uugateService *Uugate

func InitUugate(env string) {
	config := uugateDevCfg
	if uugateService == nil {
		if env == "prod" {
			config = uugateProdCfg
		}
	}
	err := json.Unmarshal([]byte(config), &uugateService)
	if err != nil {
		panic(err)
	}

	payment := Payment{
		Name:        "uugate",
		PayService:  uugateService,
		PaymentType: PayTypeUsdt,
	}

	paymentRegister(payment)
}

func UugateService() *Uugate {

	return uugateService
}

type Uugate struct {
	Uid               string `json:"uid"`
	ApiKey            string `json:"api_key"`
	Domain            string `json:"domain"`
	JumpUrl           string `json:"jump_url"`
	EffectiveDuration int    `json:"effective_duration"`
}

// UugateFd QuickNode 通用结构体/包含回调
type UugateFd struct {
	Uid       string `json:"uid"`
	Sign      string `json:"sign"`
	Timestamp string `json:"timestamp"`
	Data      string `json:"data"`
}

// UugateCallBackData 回调data
type UugateCallBackData struct {
	OrderType    string             `json:"OrderType"`
	PaymentOrder UugatePaymentOrder `json:"PaymentOrder"`
	ReceiveOrder UugateReceiveOrder
}

// UugatePaymentOrder 充值回调 uugate order
type UugatePaymentOrder struct {
	UID             string `json:"UID"`
	OrderNo         string `json:"OrderNo"`
	CustomerOrderNo string `json:"CustomerOrderNo"`
	Status          string `json:"OrderStatus"`
	FinishTime      string `json:"FinishTime"`
	Amount          string `json:"Amount"`
}

// UugateReceiveOrder 提现回调 uugate order
type UugateReceiveOrder struct {
	UID             string `json:"UID"`
	OrderNo         string `json:"OrderNo"`
	CustomerOrderNo string `json:"CustomerOrderNo"`
	Status          string `json:"OrderStatus"`
	FinishTime      string `json:"FinishTime"`
	Amount          string `json:"Amount"`
	AmountInFact    string `json:"AmountInFact"`
}

func (that *Uugate) jumpUrl() string {
	if that.JumpUrl == "" {
		return that.Domain
	}
	return that.JumpUrl
}

// uugateDepositData 请求支付数据
type uugateDepositData struct {
	Amount            string `json:"Amount"`
	Blockchain        string `json:"Blockchain"`
	CustomerOrderNo   string `json:"CustomerOrderNo"`
	EffectiveDuration int    `json:"EffectiveDuration"`
	JumpURL           string `json:"JumpURL"`
}

// uugateDepositResult 请求支付数据返回
type uugateDepositResult struct {
	CheckOutUrl    string `json:"CheckOutUrl"`    // 支付链接
	ReceiveAddress string `json:"ReceiveAddress"` // 收款地址
	Code           int    `json:"code"`           // 状态码 0 成功
	Msg            string `json:"msg"`
}

func (that *Uugate) CallDeposit(id, amount string) (CallDepositResult, error) {

	var (
		url  = that.Domain + "/Open.Customer/CreateReceiveOrder"
		resp CallDepositResult
	)

	amountDec, err := decimal.NewFromString(amount)
	if err != nil {
		fmt.Printf("UugateCallDepositErr:amount err,id:%s,amount:%s,err:%s \n", id, amount, err.Error())
		return resp, err
	}

	// 1.构造 uugate 请求体
	req := UugateFd{
		Uid:       that.Uid,
		Timestamp: fmt.Sprintf("%d", time.Now().Unix()),
	}

	data := uugateDepositData{
		Amount:            amountDec.String(),
		Blockchain:        "trc20",
		CustomerOrderNo:   id,
		EffectiveDuration: that.EffectiveDuration,
		JumpURL:           that.jumpUrl(),
	}

	dataStr, _ := json.Marshal(data)
	req.Data = string(dataStr)
	req.Sign = that.sign(req)

	payload, _ := json.Marshal(req)
	fmt.Printf("UugateCallDeposit id:%s,url:%s,req:%s\n", id, url, string(payload))

	// 2.向 uugate 发送充值请求
	respBytes, err := that.sendRequest(url, req)
	if err != nil {
		fmt.Printf("UugateCallDepositErr sendRequest id:%s,url:%s,req:%s\n", id, url, string(payload))
		return resp, err
	}

	fmt.Printf("UugateCallDeposit id:%s,url:%s,resp:%s\n", id, url, string(respBytes))

	// 3.解析uugate返回
	var payResp uugateDepositResult
	err = json.Unmarshal(respBytes, &payResp)
	if err != nil {
		fmt.Printf("UugateCallDepositErr:JsonUnmarshal err, id:%s,url:%s,req:%s resp:%s \n",
			id, url, string(payload), string(respBytes))
		return resp, err
	}

	if !(payResp.Code == 0 && payResp.Msg == "success") {
		fmt.Printf("UugateCallDepositErr:status err, id:%s,url:%s,req:%s resp:%s \n",
			id, url, string(payload), string(respBytes))
		return resp, fmt.Errorf("%s", payResp.Msg)
	}

	// 4.封装到通用返回
	resp.PayUrl = payResp.CheckOutUrl
	resp.ToAddress = payResp.ReceiveAddress

	return resp, nil
}

type uugateDepositQueryData struct {
	CustomerOrderNo string `json:"CustomerOrderNo"`
}

type uugateDepositQueryResp struct {
	ReceiveOrder UugateReceiveOrder `json:"ReceiveOrder"`
	Code         int                `json:"code"`
	Msg          string             `json:"msg"`
}

func (that *Uugate) CallDepositOrderQuery(orderId, externalOrderId string) (PaymentOrderQueryResult, error) {

	var (
		resp       PaymentOrderQueryResult
		uugateResp uugateDepositQueryResp
	)

	uugateResp, err := that.getReceiveOrderStatus(orderId)
	if err != nil {
		return resp, err
	}

	resp.ExternalStatus = uugateResp.ReceiveOrder.Status
	switch resp.ExternalStatus {
	case "待付款":
		resp.Status = OrderStatusPending
	case "付款超时":
		resp.Status = OrderStatusExpired
	case "已完成", "补单已完成", "付款风险": //付款风险是已完成的订单，付款方地址存在问题，可以视为完成
		resp.Status = OrderStatusSuccess
		// todo finishTime
	default:
		return resp, errors.New("unknown status:" + uugateResp.ReceiveOrder.Status)
	}

	resp.OrderNo = uugateResp.ReceiveOrder.CustomerOrderNo
	resp.ExternalOrderID = uugateResp.ReceiveOrder.OrderNo
	resp.Amount = uugateResp.ReceiveOrder.Amount
	resp.RealAmount = uugateResp.ReceiveOrder.AmountInFact

	return resp, nil
}

// 查询充值订单原始返回
func (that *Uugate) getReceiveOrderStatus(orderId string) (uugateDepositQueryResp, error) {

	var (
		url  = that.Domain + "/Open.Customer/GetReceiveOrderStatus"
		resp uugateDepositQueryResp
	)

	// 1.构造 uugate 请求体
	req := UugateFd{
		Uid:       that.Uid,
		Timestamp: fmt.Sprintf("%d", time.Now().UnixMilli()),
	}

	data := uugateDepositQueryData{
		CustomerOrderNo: orderId,
	}

	dataStr, _ := json.Marshal(data)
	req.Data = string(dataStr)
	req.Sign = that.sign(req)

	payload, _ := json.Marshal(req)
	fmt.Printf("UugateCallDepositOrderQuery id:%s,url:%s,req:%s\n", orderId, url, string(payload))

	// 2.向 uugate 发送请求

	bytesRes, err := that.sendRequest(url, req)
	if err != nil {
		fmt.Printf("UugateCallDepositOrderQueryErr id:%s,url:%s,req:%s\n", orderId, url, string(payload))
		return resp, err
	}
	fmt.Printf("UugateCallDepositOrderQuery id:%s,url:%s,resp:%s\n", orderId, url, string(bytesRes))

	// 3.解析返回
	err = json.Unmarshal(bytesRes, &resp)
	if err != nil {
		fmt.Printf("UugateCallDepositOrderQueryErr:JsonUnmarshal err, id:%s,url:%s,req:%s resp:%s \n",
			orderId, url, string(payload), string(bytesRes))
		return resp, err
	}

	if !(resp.Code == 0 && resp.Msg == "success") {
		fmt.Printf("UugateCallDepositOrderQueryErr:status err, id:%s,url:%s,req:%s resp:%s \n",
			orderId, url, string(payload), string(bytesRes))
		return resp, fmt.Errorf("%s", resp.Msg)
	}

	return resp, nil
}

func (that *Uugate) sign(req UugateFd) string {

	signStr := fmt.Sprintf("%s%s%s%s", that.Uid, req.Data, that.ApiKey, req.Timestamp)
	// 计算 MD5
	hash := md5.Sum([]byte(signStr))
	return hex.EncodeToString(hash[:])
}

func (that *Uugate) sendRequest(url string, params any) ([]byte, error) {

	bodyBytes, _ := json.Marshal(params)

	req, err := http.NewRequest("POST", url, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("create request error: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// 发送请求
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request error: %w", err)
	}
	defer resp.Body.Close()

	// 读取响应
	respBytes, _ := io.ReadAll(resp.Body)
	return respBytes, nil
}
