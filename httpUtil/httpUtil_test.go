package httpUtil

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// 测试Get方法
func TestGet(t *testing.T) {
	// 创建一个测试HTTP服务器
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 检查请求方法
		if r.Method != "GET" {
			t.Error("Expected GET request")
		}

		// 测试不同的路径返回不同的响应
		switch r.URL.Path {
		case "/success":
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("success response"))
		case "/error":
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("error response"))
		case "/headers":
			// 测试请求头
			if r.Header.Get("Test-Header") != "test-value" {
				t.Error("Expected Test-Header to be 'test-value'")
			}
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("headers test"))
		case "/timeout":
			// 测试超时
			time.Sleep(2 * time.Second)
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("should timeout"))
		default:
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("not found"))
		}
	}))
	defer ts.Close()

	// 创建测试客户端
	client := NewRequestClient()

	tests := []struct {
		name       string
		url        string
		setHeaders map[string]string
		setTimeout time.Duration
		wantErr    bool
		wantBody   string
	}{
		{name: "成功请求", url: ts.URL + "/success", wantErr: false, wantBody: "success response"},
		{name: "服务器错误", url: ts.URL + "/error", wantErr: false, wantBody: "error response"},
		{name: "带请求头的请求", url: ts.URL + "/headers", setHeaders: map[string]string{"Test-Header": "test-value"}, wantErr: false, wantBody: "headers test"},
		{name: "超时请求", url: ts.URL + "/timeout", setTimeout: time.Second, wantErr: true},
		{name: "无效URL", url: "http://invalid.url", wantErr: true, wantBody: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 重置客户端状态
			client = NewRequestClient()

			// 设置超时
			if tt.setTimeout > 0 {
				client.SetTimeout(tt.setTimeout)
			}

			// 设置请求头
			if tt.setHeaders != nil {
				for key, value := range tt.setHeaders {
					client.SetHeader(key, value)
				}
			}

			resp := client.Get(tt.url)

			if tt.wantErr {
				if resp.Err == nil {
					t.Error("Expected error, got nil")
				}
				return
			}

			if resp.Err != nil {
				t.Errorf("Unexpected error: %v", resp.Err)
				return
			}

			if got := resp.Text(); got != tt.wantBody {
				t.Errorf("Expected body '%s', got '%s'", tt.wantBody, got)
			}
		})
	}
}

// 测试Resp的各种方法
func TestRespMethods(t *testing.T) {
	// 测试Resp的各种方法
	jsonBody := `{"key":"value","num":123}`
	resp := &Resp{
		Body: []byte(jsonBody),
	}

	t.Run("Text", func(t *testing.T) {
		if resp.Text() != jsonBody {
			t.Error("Text() did not return expected body")
		}
	})

	t.Run("Json", func(t *testing.T) {
		jsonData, err := resp.Json()
		if err != nil {
			t.Errorf("Json() returned error: %v", err)
		}
		if jsonData["key"] != "value" || jsonData["num"] != float64(123) {
			t.Error("Json() did not return expected data")
		}
	})

	t.Run("Json with invalid data", func(t *testing.T) {
		res := &Resp{
			Body: []byte("invalid json"),
		}
		_, err := res.Json()
		if err == nil {
			t.Error("Expected error for invalid JSON, got nil")
		}
	})

	t.Run("JsonObj", func(t *testing.T) {
		type TestStruct struct {
			Key string `json:"key"`
			Num int    `json:"num"`
		}
		var obj TestStruct
		err := resp.JsonObj(&obj)
		if err != nil {
			t.Errorf("JsonObj() returned error: %v", err)
		}
		if obj.Key != "value" || obj.Num != 123 {
			t.Error("JsonObj() did not return expected data")
		}
	})
}

// 测试Post方法
func TestPost(t *testing.T) {
	// 创建一个测试用的HTTP服务器
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 检查请求头
		contentType := r.Header.Get("Content-Type")

		// 读取请求体
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Error("读取请求体失败:", err)
			return
		}

		// 根据不同的Content-Type返回不同的响应
		switch {
		case strings.Contains(contentType, "multipart/form-data"):
			w.Write([]byte("multipart response"))
		case strings.Contains(contentType, "application/json"):
			// 验证JSON数据
			if !strings.Contains(string(body), `"key1":"value1"`) {
				t.Error("JSON数据不正确")
			}
			w.Write([]byte("json response"))
		case strings.Contains(contentType, "application/x-www-form-urlencoded"):
			w.Write([]byte("form response"))
		default:
			w.Write([]byte("unknown format"))
		}
	}))
	defer ts.Close()

	// 创建测试用的临时文件
	tmpFile, err := os.CreateTemp("", "testfile")
	if err != nil {
		t.Fatal("创建临时文件失败:", err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.WriteString("test file content")
	tmpFile.Close()

	// 测试用例
	tests := []struct {
		name         string
		url          string
		data         map[string]interface{}
		setHeaders   map[string]string
		setTimeout   time.Duration
		setJson      bool
		setMultipart bool
		want         string
		wantErr      bool
	}{
		{
			name: "测试JSON格式请求",
			url:  ts.URL,
			data: map[string]interface{}{"key1": "value1", "key2": 2},
			setHeaders: map[string]string{
				"X-Custom-Header": "test",
			},
			setJson: true,
			want:    "json response",
			wantErr: false,
		},
		{
			name:    "测试表单格式请求",
			url:     ts.URL,
			data:    map[string]interface{}{"key1": "value1", "key2": 2},
			want:    "form response",
			wantErr: false,
		},
		{
			name: "测试Multipart格式请求(带文件)",
			url:  ts.URL,
			data: map[string]interface{}{
				"file": func() *os.File {
					f, _ := os.Open(tmpFile.Name())
					return f
				}(),
				"field": "value",
			},
			setMultipart: true,
			want:         "multipart response",
			wantErr:      false,
		},
		{
			name:         "测试Multipart格式请求(不带文件)",
			url:          ts.URL,
			data:         map[string]interface{}{"key1": "value1", "key2": 2},
			setMultipart: true,
			want:         "multipart response",
			wantErr:      false,
		},
		{
			name:    "测试无效URL",
			url:     "http://invalid.url",
			data:    map[string]interface{}{"key": "value"},
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建新的客户端实例
			rc := NewRequestClient()

			// 设置超时
			if tt.setTimeout > 0 {
				rc.SetTimeout(tt.setTimeout)
			}

			// 设置请求头
			if tt.setHeaders != nil {
				for key, value := range tt.setHeaders {
					rc.SetHeader(key, value)
				}
			}

			// 设置请求格式
			rc.isJson = tt.setJson
			rc.isMultipart = tt.setMultipart

			resp := rc.Post(tt.url, tt.data)

			if (resp.Err != nil) != tt.wantErr {
				t.Errorf("Post() error = %v, wantErr %v", resp.Err, tt.wantErr)
				return
			}

			if !tt.wantErr && resp.Text() != tt.want {
				t.Errorf("Post() = %v, want %v", resp.Text(), tt.want)
			}
		})
	}
}

// 测试Post方法的文件处理功能
func TestPost_FileHandling(t *testing.T) {
	// 创建一个测试用的HTTP服务器
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 解析multipart表单
		err := r.ParseMultipartForm(10 << 20) // 10MB
		if err != nil {
			t.Error("解析multipart表单失败:", err)
			return
		}

		// 检查文件是否存在
		file, header, err := r.FormFile("file")
		if err != nil {
			t.Error("获取上传文件失败:", err)
			return
		}
		defer file.Close()

		// 检查文件名
		if header.Filename != filepath.Base(header.Filename) {
			t.Error("文件名不正确")
		}

		// 检查文件内容
		buf := new(bytes.Buffer)
		buf.ReadFrom(file)
		if buf.String() != "test file content" {
			t.Error("文件内容不正确")
		}

		w.Write([]byte("file upload success"))
	}))
	defer ts.Close()

	// 创建测试用的临时文件
	tmpFile, err := os.CreateTemp("", "testfile")
	if err != nil {
		t.Fatal("创建临时文件失败:", err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.WriteString("test file content")
	tmpFile.Close()

	// 打开文件用于测试
	file, err := os.Open(tmpFile.Name())
	if err != nil {
		t.Fatal("打开临时文件失败:", err)
	}
	defer file.Close()

	rc := NewRequestClient()
	rc.isMultipart = true
	resp := rc.Post(ts.URL, map[string]interface{}{
		"file":  file,
		"field": "value",
	})

	if resp.Err != nil {
		t.Error("文件上传请求失败:", resp.Err)
	}

	if resp.Text() != "file upload success" {
		t.Error("文件上传响应不正确:", resp.Text())
	}
}

// 测试超时功能
func TestPost_Timeout(t *testing.T) {
	// 创建一个慢速响应的测试服务器
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second) // 延迟2秒响应
		w.Write([]byte("response"))
	}))
	defer ts.Close()

	rc := NewRequestClient()
	rc.SetTimeout(time.Second) // 设置1秒超时

	resp := rc.Post(ts.URL, map[string]interface{}{"key": "value"})

	if resp.Err == nil {
		t.Error("预期超时错误，但未发生")
	}
}

// 测试URL编码功能
func TestUrlEncode(t *testing.T) {
	params := map[string]interface{}{
		"key1": "value1",
		"key2": 123,
		"key3": true,
	}

	encoded := UrlEncode(params)
	expected := "key1=value1&key2=123&key3=true"
	if encoded != expected {
		t.Errorf("Expected '%s', got '%s'", expected, encoded)
	}

	// 测试空参数
	emptyParams := map[string]interface{}{}
	emptyEncoded := UrlEncode(emptyParams)
	if emptyEncoded != "" {
		t.Errorf("Expected empty string, got '%s'", emptyEncoded)
	}
}

// 测试响应头功能
func TestResponseHeaders(t *testing.T) {
	// 创建一个返回特定响应头的测试服务器
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Custom-Header", "custom-value")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	}))
	defer ts.Close()

	rc := NewRequestClient()
	resp := rc.Get(ts.URL)

	if resp.Err != nil {
		t.Errorf("Unexpected error: %v", resp.Err)
		return
	}

	if resp.Headers["X-Custom-Header"] != "custom-value" {
		t.Error("响应头X-Custom-Header不正确")
	}

	if resp.Headers["Content-Type"] != "application/json" {
		t.Error("响应头Content-Type不正确")
	}

	if resp.StatusCode != http.StatusOK {
		t.Error("状态码不正确")
	}
}
