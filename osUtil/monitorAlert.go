package osUtil

import (
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/Tomatosky/jo-util/cryptor"
	"github.com/Tomatosky/jo-util/httpUtil"
	"github.com/Tomatosky/jo-util/logger"
)

// Alert 报警接口，由用户自行实现
type Alert interface {
	Alert(resourceType string, value float64, threshold float64, duration string)
}

// defaultAlert 默认报警实现
type defaultAlert struct{}

func (d *defaultAlert) Alert(resourceType string, value float64, threshold float64, duration string) {
	logger.Log.Warn(fmt.Sprintf("[资源报警] %s 当前值: %.2f%% 阈值: %.2f%% 持续时间: %s",
		resourceType, value, threshold, duration))
}

// DingdingAlert 钉钉报警实现
type DingdingAlert struct {
	secret      string
	accessToken string
}

func NewDingdingAlert(secret string, accessToken string) *DingdingAlert {
	return &DingdingAlert{
		secret:      secret,
		accessToken: accessToken,
	}
}

func (d *DingdingAlert) Alert(resourceType string, value float64, threshold float64, duration string) {
	go func() {
		timestamp := strconv.FormatInt(time.Now().UnixMilli(), 10)
		hmacCode := cryptor.HmacSha256WithBase64(timestamp+"\n"+d.secret, d.secret)
		sign := url.QueryEscape(hmacCode)

		dingdingUrl := fmt.Sprintf("https://oapi.dingtalk.com/robot/send?access_token=%s&timestamp=%s&sign=%s", d.accessToken, timestamp, sign)
		postData := map[string]interface{}{
			"msgtype": "markdown",
			"markdown": map[string]string{
				"title": "[资源报警]",
				"text": fmt.Sprintf("### %s 当前值: %.2f%% 阈值: %.2f%% 持续时间: %s",
					resourceType, value, threshold, duration),
			},
		}

		client := httpUtil.NewRequestClient()
		client.IsJson = true
		resp := client.Post(dingdingUrl, postData)
		if resp.Err != nil {
			logger.Log.Error(fmt.Sprintf("err=%v", resp.Err))
			return
		}

		if resp.StatusCode != 200 {
			logger.Log.Error(fmt.Sprintf("请求钉钉失败,状态码:%d", resp.StatusCode))
			return
		}

		respJson, err := resp.Json()
		if err != nil {
			logger.Log.Error(fmt.Sprintf("err=%v", err))
			return
		}

		if respJson["errcode"] != 0 {
			logger.Log.Error(fmt.Sprintf("请求钉钉失败,errcode:%d,errmsg:%s", respJson["errcode"], respJson["errmsg"]))
			return
		}
	}()
}

// GotifyAlert Gotify报警实现
type GotifyAlert struct {
	host  string
	token string
}

func (g *GotifyAlert) Alert(resourceType string, value float64, threshold float64, duration string) {
	go func() {
		postData := map[string]interface{}{
			"title": "[资源报警]",
			"message": fmt.Sprintf("### %s 当前值: %.2f%% 阈值: %.2f%% 持续时间: %s",
				resourceType, value, threshold, duration),
			"extras": map[string]any{
				"client::display": map[string]string{
					"contentType": "text/markdown",
				},
			},
		}

		gotifyUrl := fmt.Sprintf("%s/message?token=%s", g.host, g.token)
		client := httpUtil.NewRequestClient()
		client.IsJson = true
		resp := client.Post(gotifyUrl, postData)
		if resp.Err != nil {
			logger.Log.Error(fmt.Sprintf("err=%v", resp.Err))
			return
		}

		if resp.StatusCode != 200 {
			logger.Log.Error(fmt.Sprintf("请求Gotify失败,状态码:%d", resp.StatusCode))
			return
		}
	}()
}
