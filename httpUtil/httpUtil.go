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
	"runtime/debug"
	"strings"
	"time"

	"github.com/Tomatosky/jo-util/convertor"
)

// RequestClient 是一个HTTP客户端工具类
type RequestClient struct {
	client      *http.Client
	headers     map[string]string
	isJson      bool
	isMultipart bool
}

// NewRequestClient 创建一个新的RequestClient实例
func NewRequestClient() *RequestClient {
	return &RequestClient{
		client: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			},
		},
		headers: make(map[string]string),
	}
}

func (rc *RequestClient) SetProxy(proxy string) {
	proxyUrl, err := url.Parse(proxy)
	if err != nil {
		debug.PrintStack()
		panic(err)
	}
	rc.client.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
		Proxy: http.ProxyURL(proxyUrl),
	}
}

func (rc *RequestClient) SetTimeout(timeout time.Duration) {
	rc.client.Timeout = timeout
}

func (rc *RequestClient) SetHeader(key, value string) {
	rc.headers[key] = value
}

func (rc *RequestClient) SetJson(isJson bool) {
	rc.isJson = isJson
}

func (rc *RequestClient) SetMultipart(isMultipart bool) {
	rc.isMultipart = isMultipart
}

type Resp struct {
	Err        error
	Body       []byte
	StatusCode int
	Headers    map[string]string
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
func (rc *RequestClient) Get(url string) *Resp {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return &Resp{Err: err}
	}

	// 设置请求头
	if len(rc.headers) > 0 {
		for key, value := range rc.headers {
			req.Header.Set(key, value)
		}
	}

	// 发送请求
	resp, err := rc.client.Do(req)
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
	case rc.isMultipart:
		// multipart/form-data 格式
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)

		for key, value := range data {
			// 处理文件上传情况
			if file, ok := value.(*os.File); ok {
				part, err := writer.CreateFormFile(key, filepath.Base(file.Name()))
				if err != nil {
					return &Resp{Err: err}
				}
				_, err = io.Copy(part, file)
				if err != nil {
					return &Resp{Err: err}
				}
				continue
			}

			// 处理普通字段
			strValue := fmt.Sprintf("%v", value)
			err := writer.WriteField(key, strValue)
			if err != nil {
				return &Resp{Err: err}
			}
		}

		err := writer.Close()
		if err != nil {
			return &Resp{Err: err}
		}

		reqBody = body
		contentType = writer.FormDataContentType()

	case rc.isJson:
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
	if len(rc.headers) > 0 {
		for key, value := range rc.headers {
			req.Header.Set(key, value)
		}
	}
	// 发送请求
	resp, err := rc.client.Do(req)
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
