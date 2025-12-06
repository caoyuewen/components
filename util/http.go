package util

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"time"
)

func Post(uri string, params map[string]string) ([]byte, error) {
	values := url.Values{}

	for k, v := range params {
		values.Add(k, v)
	}

	resp, err := http.PostForm(uri, values)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	return body, err
}

func PostMultipartFormData(uri string, params map[string]string) ([]byte, error) {
	formBuf := new(bytes.Buffer)
	writer := multipart.NewWriter(formBuf)
	defer writer.Close()
	// 写入请求参数
	for k, v := range params {
		writer.WriteField(k, v)
	}
	// 创建请求
	req, err := http.NewRequest("POST", uri, formBuf)
	if err != nil {
		return nil, err
	}
	// 设置 Content-Type
	contentType := writer.FormDataContentType()
	req.Header.Set("Content-Type", contentType)

	// 同步执行请求
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	// 读取应答报文
	body, err := io.ReadAll(resp.Body)
	return body, err
}

func Get(url string) ([]byte, error) {
	return GetWithTimeout(url, 30*time.Second)
}

func GetWithTimeout(url string, timeout time.Duration) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	request, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	return io.ReadAll(response.Body)
}

func PostTpl[T any](url string, req any) (*T, error) {
	return PostTplWithTimeout[T](url, req, 30*time.Second)
}

func PostTplWithTimeout[T any](url string, req any, timeout time.Duration) (*T, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	data, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	reqInst, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	reqInst.Header.Set("Content-Type", "application/json")
	respInst, err := http.DefaultClient.Do(reqInst)
	if err != nil {
		return nil, err
	}
	defer respInst.Body.Close()

	body, err := io.ReadAll(respInst.Body)
	if err != nil {
		return nil, err
	}

	res := new(T)
	if err = json.Unmarshal(body, res); err != nil {
		return nil, err
	}

	return res, nil
}
