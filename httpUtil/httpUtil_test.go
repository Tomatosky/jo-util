package httpUtil

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewRequestClient(t *testing.T) {
	client := NewRequestClient()
	if client.Client == nil {
		t.Error("Expected http.Client to be initialized, got nil")
	}
}

func TestDownload(t *testing.T) {
	// 创建测试服务器
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("test response"))
	}))
	defer ts.Close()

	client := NewRequestClient()

	// 测试正常下载
	body, err := client.Download(ts.URL, &GetOptions{})
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if string(body) != "test response" {
		t.Errorf("Expected 'test response', got '%s'", string(body))
	}

	// 测试带header的下载
	headers := map[string]string{"X-Test": "value"}
	body, err = client.Download(ts.URL, &GetOptions{Headers: headers})
	if err != nil {
		t.Errorf("Expected no error with headers, got %v", err)
	}

	// 测试超时
	_, err = client.Download(ts.URL, &GetOptions{Timeout: 1})
	if err != nil {
		t.Errorf("Expected no error with timeout, got %v", err)
	}

	// 测试无效URL
	_, err = client.Download("invalid-url", &GetOptions{})
	if err == nil {
		t.Error("Expected error with invalid URL, got nil")
	}
}

func TestGet(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("test get response"))
	}))
	defer ts.Close()

	client := NewRequestClient()

	// 测试正常GET请求
	response, err := client.Get(ts.URL, &GetOptions{})
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if response != "test get response" {
		t.Errorf("Expected 'test get response', got '%s'", response)
	}
}

func TestPost(t *testing.T) {
	// 测试表单POST
	formServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != "application/x-www-form-urlencoded" {
			t.Error("Expected Content-Type to be application/x-www-form-urlencoded")
		}
		w.Write([]byte("form post response"))
	}))
	defer formServer.Close()

	// 测试JSON POST
	jsonServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != "application/json" {
			t.Error("Expected Content-Type to be application/json")
		}

		var data map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			t.Errorf("Error decoding JSON: %v", err)
		}
		if data["key"] != "value" {
			t.Errorf("Expected key=value, got %v", data)
		}
		w.Write([]byte("json post response"))
	}))
	defer jsonServer.Close()

	client := NewRequestClient()

	// 测试表单POST
	data := map[string]interface{}{"key": "value"}
	response, err := client.Post(formServer.URL, data, &PostOptions{IsJson: false})
	if err != nil {
		t.Errorf("Expected no error with form post, got %v", err)
	}
	if response != "form post response" {
		t.Errorf("Expected 'form post response', got '%s'", response)
	}

	// 测试JSON POST
	response, err = client.Post(jsonServer.URL, data, &PostOptions{IsJson: true})
	if err != nil {
		t.Errorf("Expected no error with json post, got %v", err)
	}
	if response != "json post response" {
		t.Errorf("Expected 'json post response', got '%s'", response)
	}

	// 测试带header的POST
	headers := map[string]string{"X-Test": "value"}
	response, err = client.Post(formServer.URL, data, &PostOptions{Headers: headers})
	if err != nil {
		t.Errorf("Expected no error with headers, got %v", err)
	}

	// 测试超时
	response, err = client.Post(formServer.URL, data, &PostOptions{Timeout: 1})
	if err != nil {
		t.Errorf("Expected no error with timeout, got %v", err)
	}

	// 测试无效URL
	_, err = client.Post("invalid-url", data, &PostOptions{})
	if err == nil {
		t.Error("Expected error with invalid URL, got nil")
	}
}

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
}
