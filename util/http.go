package util

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
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
	client := &http.Client{}

	//提交请求
	request, err := http.NewRequest("GET", url, nil)
	//异常捕捉
	if err != nil {
		panic(err)
	}

	//处理返回结果
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	//关闭流
	defer response.Body.Close()
	//检出结果集
	resp, err := io.ReadAll(response.Body)

	return resp, err
}

func PostTpl[T any](url string, req any) (*T, error) {

	data, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	reqInst, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	reqInst.Header.Set("Content-Type", "application/json")
	respInst, err := http.DefaultClient.Do(reqInst)
	if err != nil {
		return nil, err
	}

	defer respInst.Body.Close()
	body, err := ioutil.ReadAll(respInst.Body)
	if err != nil {
		return nil, err
	}

	res := new(T)
	if err = json.Unmarshal(body, res); err != nil {
		return nil, err
	}

	return res, nil
}
