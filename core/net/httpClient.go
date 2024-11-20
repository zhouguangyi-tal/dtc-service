package net

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type ResponseType[T any] struct {
	Code   int    `json:"code"`
	Msg    string `json:"msg"`
	Result T      `json:"result"`
}

// GetRequest 发送不带参数的GET请求并返回解析后的响应体
func GetRequest[T any](url string) (*T, error) {
	// 创建一个HTTP客户端
	client := &http.Client{}

	// 创建一个GET请求
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request failed with status code: %d", resp.StatusCode)
	}

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// 解析JSON响应体
	var result T
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetRequestWithParams 发送带有查询参数的GET请求并返回解析后的响应体
func GetRequestWithParams[T any](baseURL string, params map[string]string) (*ResponseType[T], error) {
	// 解析基本URL
	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}

	// 添加查询参数
	q := u.Query()
	for key, value := range params {
		q.Set(key, value)
	}
	u.RawQuery = q.Encode()

	// 创建一个HTTP客户端
	client := &http.Client{}

	// 创建一个GET请求
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request failed with status code: %d", resp.StatusCode)
	}

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// 解析JSON响应体
	var res ResponseType[any]
	if err := json.Unmarshal(body, &res); err != nil {
		return nil, err
	}
	if res.Code != 0 {
		return nil, fmt.Errorf("request failed with status code: %", res.Code, res.Msg)
	}

	var result ResponseType[T]
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// PostRequest 发送POST请求并返回解析后的响应体
func PostRequest[T any](baseURL string, params map[string]string, payload any) (*ResponseType[T], error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}

	// 添加查询参数
	q := u.Query()
	for key, value := range params {
		q.Set(key, value)
	}
	u.RawQuery = q.Encode()

	// 将payload编码为JSON
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	// 创建一个HTTP客户端
	client := &http.Client{}

	// 创建POST请求
	req, err := http.NewRequest("POST", u.String(), bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request failed with status code: %d", resp.StatusCode)
	}

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	// 解析JSON响应体
	var res ResponseType[any]
	if err := json.Unmarshal(body, &res); err != nil {
		return nil, err
	}
	if res.Code != 0 {
		return nil, fmt.Errorf("request failed with status code: %", res.Code, res.Msg)
	}

	var result ResponseType[T]
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	return &result, nil
}
