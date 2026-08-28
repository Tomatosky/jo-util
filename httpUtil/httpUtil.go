package httpUtil

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Tomatosky/jo-util/convertor"
	"github.com/Tomatosky/jo-util/logger"
)

// RequestClient 是一个HTTP客户端工具类
type RequestClient struct {
	Client      *http.Client
	Headers     map[string]string
	IsJson      bool
	IsMultipart bool
}

// NewRequestClient 创建一个新的RequestClient实例
func NewRequestClient() *RequestClient {
	return &RequestClient{
		Client: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			},
		},
		Headers: make(map[string]string),
	}
}

func (rc *RequestClient) SetProxy(proxy string) {
	proxyUrl, err := url.Parse(proxy)
	if err != nil {
		logger.Log.Error(fmt.Sprintf("%v", err))
		panic(err)
	}
	rc.Client.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
		Proxy: http.ProxyURL(proxyUrl),
	}
}

func (rc *RequestClient) SetTimeout(timeout time.Duration) {
	rc.Client.Timeout = timeout
}

type Resp struct {
	Err        error
	Body       []byte
	StatusCode int
	Headers    map[string]string
}

func buildMultipartBody(data map[string]interface{}) (body *bytes.Buffer, contentType string, err error) {
	body = &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	contentType = writer.FormDataContentType()
	defer func() {
		if closeErr := writer.Close(); err == nil {
			err = closeErr
		}
	}()

	for key, value := range data {
		if file, ok := value.(*os.File); ok {
			part, createErr := writer.CreateFormFile(key, filepath.Base(file.Name()))
			if createErr != nil {
				return nil, "", createErr
			}
			if _, err = io.Copy(part, file); err != nil {
				return nil, "", err
			}
			continue
		}

		if err = writer.WriteField(key, fmt.Sprintf("%v", value)); err != nil {
			return nil, "", err
		}
	}

	return body, contentType, nil
}

func (r *Resp) Text() string {
	return string(r.Body)
}

func (r *Resp) Json() (map[string]any, error) {
	jsonData := map[string]any{}
	err := json.Unmarshal(r.Body, &jsonData)
	return jsonData, err
}

func (r *Resp) JsonObj(v any) error {
	err := json.Unmarshal(r.Body, v)
	return err
}

// Get 发送GET请求
func (rc *RequestClient) Get(requestURL string) *Resp {
	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		return &Resp{Err: err}
	}

	// 设置请求头
	if len(rc.Headers) > 0 {
		for key, value := range rc.Headers {
			req.Header.Set(key, value)
		}
	}

	// 发送请求
	resp, err := rc.Client.Do(req)
	if err != nil {
		return &Resp{Err: err}
	}
	defer resp.Body.Close()

	// 读取响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return &Resp{Err: err}
	}

	// 收集响应头
	headers := make(map[string]string)
	for k, v := range resp.Header {
		if len(v) > 0 {
			headers[k] = v[0] // 只取第一个值
		}
	}

	return &Resp{
		Err:        nil,
		Body:       body,
		StatusCode: resp.StatusCode,
		Headers:    headers,
	}
}

// Post 发送POST请求
func (rc *RequestClient) Post(postUrl string, data map[string]interface{}) *Resp {
	var (
		reqBody     io.Reader
		contentType string
	)
	// 判断请求体格式
	switch {
	case rc.IsMultipart:
		// multipart/form-data 格式
		body, multipartContentType, err := buildMultipartBody(data)
		if err != nil {
			return &Resp{Err: err}
		}
		reqBody = body
		contentType = multipartContentType

	case rc.IsJson:
		// JSON 格式
		jsonData, err := json.Marshal(data)
		if err != nil {
			return &Resp{Err: err}
		}
		reqBody = bytes.NewBuffer(jsonData)
		contentType = "application/json"

	default:
		// 表单格式 (x-www-form-urlencoded)
		formData := url.Values{}
		for key, value := range data {
			formData.Set(key, fmt.Sprintf("%v", value)) // 确保值转为字符串
		}
		reqBody = strings.NewReader(formData.Encode())
		contentType = "application/x-www-form-urlencoded"
	}

	// 创建 HTTP 请求
	req, err := http.NewRequest("POST", postUrl, reqBody)
	if err != nil {
		return &Resp{Err: err}
	}
	// 设置请求头
	req.Header.Set("Content-Type", contentType)
	if len(rc.Headers) > 0 {
		for key, value := range rc.Headers {
			req.Header.Set(key, value)
		}
	}
	// 发送请求
	resp, err := rc.Client.Do(req)
	if err != nil {
		return &Resp{Err: err}
	}
	defer resp.Body.Close()
	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return &Resp{Err: err}
	}

	// 收集响应头
	headers := make(map[string]string)
	for k, v := range resp.Header {
		if len(v) > 0 {
			headers[k] = v[0] // 只取第一个值
		}
	}

	return &Resp{
		Err:        nil,
		Body:       body,
		StatusCode: resp.StatusCode,
		Headers:    headers,
	}
}

func UrlEncode(params map[string]interface{}) string {
	var p = url.Values{}
	for k, v := range params {
		p.Add(k, convertor.ToString(v))
	}
	return p.Encode()
}
