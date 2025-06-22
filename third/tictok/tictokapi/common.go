package tictokapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/caoyuewen/components/util"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
)

const (
	ApiCodeSuccess = 0
)

func Post(url string, body interface{}, header map[string]string) ([]byte, error) {
	// 构造请求体
	log.Debug("-------------------------- tictok post --------------------------")
	log.Debug("Post url:", url)
	log.Debug("Post header:", util.ToJson(header))
	log.Debug("Post body:", util.ToJson(body))
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %v", err)
	}

	// 创建HTTP请求
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	// 设置请求头
	for k, v := range header {
		req.Header.Set(k, v)
	}

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// 读取响应内容
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}
	log.Debugf("Post resp:%s", string(bodyBytes))
	return bodyBytes, nil
}
