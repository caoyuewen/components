package pay

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// ==================== TRON 网络配置 ====================

const (
	TronGridAPI = "https://api.trongrid.io" // TronGrid API (主网)

	UsdtContractAddress = "TR7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t" // USDT TRC20 合约地址

	UsdtDecimals = 6 // USDT 精度 (6位小数)
)

var (
	tronAPIKey string // TronGrid API Key
)

// SetTronAPIKey 设置 TronGrid API Key
func SetTronAPIKey(apiKey string) {
	tronAPIKey = apiKey
}

// ==================== TRON API 响应结构 ====================

// TRC20Transaction TRC20 交易记录
type TRC20Transaction struct {
	TransactionID  string    `json:"transaction_id"`
	TokenInfo      TokenInfo `json:"token_info"`
	BlockTimestamp int64     `json:"block_timestamp"`
	From           string    `json:"from"`
	To             string    `json:"to"`
	Type           string    `json:"type"`
	Value          string    `json:"value"`
}

// TokenInfo 代币信息
type TokenInfo struct {
	Symbol   string `json:"symbol"`
	Address  string `json:"address"`
	Decimals int    `json:"decimals"`
	Name     string `json:"name"`
}

// TRC20Response TRC20 交易查询响应
type TRC20Response struct {
	Data    []TRC20Transaction `json:"data"`
	Success bool               `json:"success"`
	Meta    struct {
		At          int64  `json:"at"`
		Fingerprint string `json:"fingerprint"`
		PageSize    int    `json:"page_size"`
	} `json:"meta"`
}

// ==================== TRON API 方法 ====================

// GetTRC20Transactions 获取地址的 TRC20 交易记录
func GetTRC20Transactions(address string, minTimestamp int64, limit int) ([]TRC20Transaction, error) {
	if limit <= 0 {
		limit = 50
	}

	url := fmt.Sprintf("%s/v1/accounts/%s/transactions/trc20?only_to=true&limit=%d&contract_address=%s",
		TronGridAPI, address, limit, UsdtContractAddress)

	if minTimestamp > 0 {
		url += fmt.Sprintf("&min_timestamp=%d", minTimestamp)
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")
	if tronAPIKey != "" {
		req.Header.Set("TRON-PRO-API-KEY", tronAPIKey)
	}

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result TRC20Response
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("parse response error: %v, body: %s", err, string(body))
	}

	if !result.Success {
		return nil, fmt.Errorf("trongrid api error: %s", string(body))
	}

	return result.Data, nil
}

// GetTransactionInfo 获取交易详情
func GetTransactionInfo(txHash string) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/v1/transactions/%s", TronGridAPI, txHash)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")
	if tronAPIKey != "" {
		req.Header.Set("TRON-PRO-API-KEY", tronAPIKey)
	}

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	return result, nil
}
