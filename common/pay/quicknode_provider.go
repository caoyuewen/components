package pay

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"strings"
	"time"

	"github.com/btcsuite/btcutil/base58"
	"github.com/caoyuewen/components/common/caches"
	"github.com/shopspring/decimal"
	log "github.com/sirupsen/logrus"
)

const (
	network               = "tron-mainnet"
	QuickNodeWebhooksName = "tron usdt node webhook"
)

type QuickNodeFactory struct{}

var prodCfg = `{"api_key":"QN_d0b77e2d5daf4ada8030bf658b602f72","notify_email":"liufengzhx@gmail.com","callback":"https://effic.in/payment/callback/quicknode","jump_url":"https://effic.in","Domain":"https://snowy-divine-bridge.tron-mainnet.quiknode.pro/38651257af9f5ef32ca03dce1a09994b1003d1fb/jsonrpc"}`
var devCfg = `{"api_key":"QN_d0b77e2d5daf4ada8030bf658b602f72","notify_email":"liufengzhx@gmail.com","callback":"https://effic.in/payment/callback/quicknode","jump_url":"https://effic.in","Domain":"https://snowy-divine-bridge.tron-mainnet.quiknode.pro/38651257af9f5ef32ca03dce1a09994b1003d1fb/jsonrpc"}`

type QuickNode struct {
	ApiKey      string `json:"api_key"`
	NotifyEmail string `json:"notify_email"`
	Callback    string `json:"callback"`
	JumpUrl     string `json:"jump_url"`
	Domain      string `json:"domain"`
}

var quickNodeService *QuickNode

func InitQuickNode(env string) {
	config := devCfg
	if quickNodeService == nil {
		if env == "prod" {
			config = prodCfg
		}
	}
	err := json.Unmarshal([]byte(config), &quickNodeService)
	if err != nil {
		panic(err)
	}
}

func QuickNodeService() *QuickNode {

	return quickNodeService
}

func (that *QuickNode) CallDeposit(id, amount string) (CallDepositResult, error) {

	var res CallDepositResult

	amountDec, err := decimal.NewFromString(amount)
	if err != nil {
		log.Errorf("QuickNodeCallDepositErr:amount err,id:%s,amount:%s,err:%s \n", id, amount, err.Error())
		return res, err
	}

	address, err := caches.UsdtAddress.Pop(amountDec)
	if err != nil {
		log.Errorf("QuickNodeCallDepositErr:UsdtAddressPop err,id:%s,amount:%s,err:%s \n", id, amount, err.Error())
		return res, err
	}

	res.ToAddress = address

	return res, nil
}

func (that *QuickNode) CallDepositOrderQuery(orderId, externalOrderId string) (PaymentOrderQueryResult, error) {

	return PaymentOrderQueryResult{}, nil
}

// CheckWebhooksConfig 检查并更新配置
func (that *QuickNode) CheckWebhooksConfig(wallets []string) error {

	// 1.获取webhooks列表
	list, err := that.WebhooksList()
	if err != nil {
		log.Error("QuickNodeCheckWebhooksConfig:WebhooksList err:", err)
		return err
	}

	// 2.遍历列表 是否存在匹配的webhooks名称
	for _, v := range list.Data {
		if v.Name == QuickNodeWebhooksName { // 如果存在则删除重建
			err := that.WebhooksDelete(v.Id)
			if err != nil {
				log.Error("QuickNodeCheckWebhooksConfig:WebhooksDelete err:", err)
				return err
			}
		}
	}

	// 如果地址池为空则无需创建
	if len(wallets) <= 0 {
		log.Error("Warning CheckWebhooksConfig wallets is empty")
		return nil
	}

	// 3.如果不存在则创建
	_, err = that.CreateWebhook(QuickNodeWebhooksName, wallets)
	if err != nil {
		log.Error("QuickNodeCheckWebhooksConfig:CreateWebhook err:", err)
		return err
	}
	return nil

}

// CreateWebhook 创建一个新的 webhook
func (that *QuickNode) CreateWebhook(name string, wallets []string) ([]byte, error) {

	url := "https://api.quicknode.com/webhooks/rest/v1/webhooks/template/evmWalletFilter"
	apiKey := that.ApiKey

	var evWallets []string

	// tron 地址需要转化成 EVM address
	for _, v := range wallets {
		ew, _ := TronToEvmAddress(v)
		evWallets = append(evWallets, ew)
	}

	fmt.Println("wallets", evWallets)
	payload := map[string]interface{}{
		"name":               name,
		"network":            network,
		"notification_email": that.NotifyEmail,
		"destination_attributes": map[string]interface{}{
			"url":         that.Callback,
			"compression": "none",
		},
		"status": "active",
		"templateArgs": map[string]interface{}{
			"contractAddress": "TXLAQ63Xg1NAzckPwKHvzw7CSEmLMEqcdj",
			"wallets":         evWallets,
		},
	}

	header := map[string]string{
		"accept":       "application/json",
		"Content-Type": "application/json",
		"x-api-key":    apiKey,
	}

	return that.sendRequest(url, "POST", header, payload)
}

type WebHooksListResult struct {
	Data     []WebHooksListData   `json:"data"`
	PageInfo WebHooksListPageInfo `json:"pageInfo"`
}

// WebhooksList 查询 webhooks 列表
func (that *QuickNode) WebhooksList() (WebHooksListResult, error) {

	url := "https://api.quicknode.com/webhooks/rest/v1/webhooks"
	apiKey := that.ApiKey

	header := map[string]string{
		"accept":    "application/json",
		"x-api-key": apiKey,
	}

	var res WebHooksListResult
	respBytes, err := that.sendRequest(url, "GET", header, nil)
	if err != nil {
		return res, err
	}

	err = json.Unmarshal(respBytes, &res)
	if err != nil {
		return res, err
	}

	return res, nil
}

// WebhookUpdate 只适合更新 状态 回调地址 email
func (that *QuickNode) WebhookUpdate(id string, status string) ([]byte, error) {

	url := fmt.Sprintf("https://api.quicknode.com/webhooks/rest/v1/webhooks/%s", id)

	payload := map[string]interface{}{
		"name":               QuickNodeWebhooksName,
		"network":            network,
		"notification_email": that.NotifyEmail,
		"destination_attributes": map[string]interface{}{
			"url":         that.Callback,
			"compression": "none",
		},
		"status": status,
	}

	header := map[string]string{
		"accept":       "application/json",
		"Content-Type": "application/json",
		"x-api-key":    that.ApiKey,
	}

	return that.sendRequest(url, "PATCH", header, payload)
}

// WebhooksDelete 删除
func (that *QuickNode) WebhooksDelete(id string) error {

	url := fmt.Sprintf("https://api.quicknode.com/webhooks/rest/v1/webhooks/%s", id)

	header := map[string]string{
		"accept":       "application/json",
		"Content-Type": "application/json",
		"x-api-key":    that.ApiKey,
	}

	_, err := that.sendRequest(url, "DELETE", header, nil)

	return err

}

// GetTxDetailByHash 根据交易哈希获取完整信息（包括交易金额、状态、TRON 地址）
func (that *QuickNode) GetTxDetailByHash(txHash string) (TransferInfoData, error) {

	receipt, err := that.EthGetTransactionReceipt(txHash)
	if err != nil {
		return TransferInfoData{}, err
	}

	info := TransferLogInfo(receipt.Result.Logs)

	// 交易状态
	status := "failed"
	if receipt.Result.Status == "0x1" {
		status = "success"
	}

	res := TransferInfoData{
		TxHash:    receipt.Result.TransactionHash,
		Amount:    info.Amount,
		From:      info.From, // 事件里的 from (资金来源)
		To:        info.To,   // 事件里的 to (资金接收方)
		Status:    status,
		OrgStatus: receipt.Result.Status,
	}

	return res, nil
}

type TransferLogData struct {
	From     string          `json:"from"`     // 发送放tron地址
	To       string          `json:"to"`       // 接收方tron地址
	Amount   decimal.Decimal `json:"amount"`   // 交易金额
	Contract string          `json:"contract"` // 合约地址
}

func TransferLogInfo(logs []TransferLog) TransferLogData {

	var (
		fromHex  string
		toHex    string
		amount   decimal.Decimal
		Contract string
	)

	for _, log := range logs {
		if len(log.Topics) >= 3 && log.Topics[0] == "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef" {
			// 事件里的 from/to 地址 (EVM hex 地址)
			Contract = log.Topics[0]
			fromHex = "0x" + log.Topics[1][26:]
			toHex = "0x" + log.Topics[2][26:]
			valueInt, _ := new(big.Int).SetString(log.Data[2:], 16)
			amount = decimal.NewFromBigInt(valueInt, -6) // -6 精度就是除以 1e6

			break
		}
	}

	// 事件里的 from/to 转换为 TRON 地址
	fromTron, _ := EvmToTronAddress(fromHex)
	toTron, _ := EvmToTronAddress(toHex)

	res := TransferLogData{
		From:     fromTron,
		To:       toTron,
		Amount:   amount,
		Contract: Contract,
	}

	return res
}

type EthGetTransactionReceiptResult struct {
	Result struct {
		BlockHash         string        `json:"blockHash"`
		BlockNumber       string        `json:"blockNumber"`
		ContractAddress   interface{}   `json:"contractAddress"`
		CumulativeGasUsed string        `json:"cumulativeGasUsed"`
		EffectiveGasPrice string        `json:"effectiveGasPrice"`
		From              string        `json:"from"`
		GasUsed           string        `json:"gasUsed"`
		Logs              []TransferLog `json:"logs"`
		LogsBloom         string        `json:"logsBloom"`
		Status            string        `json:"status"`
		To                string        `json:"to"`
		TransactionHash   string        `json:"transactionHash"`
		TransactionIndex  string        `json:"transactionIndex"`
		Type              string        `json:"type"`
	} `json:"result"`
}

func (that *QuickNode) EthGetTransactionReceipt(txHash string) (EthGetTransactionReceiptResult, error) {

	var result EthGetTransactionReceiptResult

	header := map[string]string{"Content-Type": "application/json"}

	req := map[string]interface{}{
		"jsonrpc": "2.0", "id": time.Now().UnixNano(),
		"method": "eth_getTransactionReceipt",
		"params": []interface{}{txHash},
	}

	resp, err := that.sendRequest(that.Domain, "POST", header, req)
	if err != nil {
		return result, err
	}

	if err := json.Unmarshal(resp, &result); err != nil {
		return result, err
	}

	return result, nil
}

// TransferInfoData 交易信息
type TransferInfoData struct {
	TxHash    string          // 交易哈希
	Amount    decimal.Decimal // 转账金额 (USDT 6位精度)
	From      string          // Transfer 事件里的 from (资金转出方)
	To        string          // Transfer 事件里的 to (资金接收方)
	Status    string          // success / failed
	OrgStatus string          // 原始状态 (0x1 / 0x0)
	//TxFrom    string          // 交易发起方 (谁发起的交易)
	//TxTo      string          // 交易目标 (通常是合约地址)
	//Gas       string          // Gas 消耗
}

// TransferInfo 根据 quickNode 回调 获取交易信息
func TransferInfo(payload []byte) (TransferInfoData, error) {

	var cb QuickNodeCallbackFd
	err := json.Unmarshal(payload, &cb)
	if err != nil {
		return TransferInfoData{}, err
	}

	if len(cb.MatchingReceipts) == 0 || len(cb.MatchingReceipts[0].Logs) == 0 {
		return TransferInfoData{}, errors.New("MatchingReceipts is empty")
	}

	receipt := cb.MatchingReceipts[0]

	info := TransferLogInfo(receipt.Logs)

	// 交易状态
	status := "failed"
	if receipt.Status == "0x1" {
		status = "success"
	}

	res := TransferInfoData{
		TxHash:    NormalizeTxHash(receipt.TransactionHash),
		Amount:    info.Amount,
		From:      info.From,
		To:        info.To,
		Status:    status,
		OrgStatus: receipt.Status,
	}

	return res, nil
}

// NormalizeTxHash 去掉前缀 0x 并转小写
func NormalizeTxHash(hash string) string {
	if len(hash) > 2 && (hash[0:2] == "0x" || hash[0:2] == "0X") {
		return strings.ToLower(hash[2:])
	}
	return strings.ToLower(hash)
}

// sendRequest 统一的请求
func (that *QuickNode) sendRequest(url, method string, header map[string]string, param map[string]any) ([]byte, error) {
	log.Info("QuickNode request url:", url)

	var reqBody io.Reader
	if param != nil {
		payload, err := json.Marshal(param)
		if err != nil {
			return nil, fmt.Errorf("json marshal error: %v", err)
		}
		log.Info("QuickNode request params:", string(payload))
		reqBody = bytes.NewReader(payload)
	}

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("create request error: %v", err)
	}

	for k, v := range header {
		req.Header.Set(k, v)
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		log.Info("QuickNode request err:", err.Error())
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response error: %v", err)
	}

	log.Info("QuickNode resp:", string(body))

	if resp.StatusCode != http.StatusOK &&
		resp.StatusCode != http.StatusCreated && // 201
		resp.StatusCode != http.StatusNoContent { // 204
		log.Info("QuickNode request status err:", resp.StatusCode)
		return nil, fmt.Errorf("status code is %d", resp.StatusCode)
	}

	return body, nil
}

// TronToEvmAddress Tron → EVM 0x 地址
func TronToEvmAddress(tronAddr string) (string, error) {
	decoded := base58.Decode(tronAddr)

	if len(decoded) < 21 {
		return "", errors.New("invalid tron address length")
	}

	// Remove the prefix (0x41) and convert to hex
	evmAddrBytes := decoded[1:21]
	return "0x" + hex.EncodeToString(evmAddrBytes), nil
}

// EvmToTronAddress EVM 0x 地址 → Tron T 开头地址
func EvmToTronAddress(evmAddr string) (string, error) {
	evmAddr = strings.ToLower(strings.TrimPrefix(evmAddr, "0x"))
	if len(evmAddr) != 40 {
		return "", fmt.Errorf("invalid evm address length")
	}
	addrBytes, err := hex.DecodeString(evmAddr)
	if err != nil {
		return "", err
	}

	// Tron 地址前缀 0x41
	tronAddr := append([]byte{0x41}, addrBytes...)

	// 计算 double-SHA256 校验和
	first := sha256.Sum256(tronAddr)
	second := sha256.Sum256(first[:])
	checksum := second[:4]

	// 拼接 + Base58 编码
	full := append(tronAddr, checksum...)
	return base58.Encode(full), nil
}

// ------------------ Response Data  Structs ------------------

type WebHooksListData struct {
	Id                    string                `json:"id"`
	Name                  string                `json:"name"`
	Status                string                `json:"status"`
	CreatedAt             time.Time             `json:"created_at"`
	DestinationAttributes DestinationAttributes `json:"destination_attributes"`
	FilterFunction        string                `json:"filter_function"`
	Network               string                `json:"network"`
	NotificationEmail     string                `json:"notification_email"`
	Sequence              int                   `json:"sequence"`
	UpdatedAt             time.Time             `json:"updated_at"`
}

type DestinationAttributes struct {
	Url           string `json:"url"`
	SecurityToken string `json:"security_token"`
	Compression   string `json:"compression"`
}

type WebHooksListPageInfo struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
	Total  int `json:"total"`
}

type QuickNodeCallbackFd struct {
	MatchingReceipts []struct {
		BlockHash         string        `json:"blockHash"`
		BlockNumber       string        `json:"blockNumber"`
		ContractAddress   interface{}   `json:"contractAddress"`
		CumulativeGasUsed string        `json:"cumulativeGasUsed"`
		EffectiveGasPrice string        `json:"effectiveGasPrice"`
		From              string        `json:"from"`
		GasUsed           string        `json:"gasUsed"`
		Logs              []TransferLog `json:"logs"`
		LogsBloom         string        `json:"logsBloom"`
		Status            string        `json:"status"`
		To                string        `json:"to"`
		TransactionHash   string        `json:"transactionHash"`
		TransactionIndex  string        `json:"transactionIndex"`
		Type              string        `json:"type"`
	} `json:"matchingReceipts"`
	MatchingTransactions interface{} `json:"matchingTransactions"`
}

type TransferLog struct {
	Address          string   `json:"address"`
	BlockHash        string   `json:"blockHash"`
	BlockNumber      string   `json:"blockNumber"`
	Data             string   `json:"data"`
	LogIndex         string   `json:"logIndex"`
	Removed          bool     `json:"removed"`
	Topics           []string `json:"topics"`
	TransactionHash  string   `json:"transactionHash"`
	TransactionIndex string   `json:"transactionIndex"`
}
