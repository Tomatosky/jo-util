package alertUtil

import (
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/Tomatosky/jo-util/cryptor"
	"github.com/Tomatosky/jo-util/httpUtil"
	"github.com/Tomatosky/jo-util/logger"
)

// DingdingAlert 钉钉报警实现
type DingdingAlert struct {
	secret      string
	accessToken string
}

func (d *DingdingAlert) Alert(title string, content string) {
	go func() {
		timestamp := strconv.FormatInt(time.Now().UnixMilli(), 10)
		hmacCode := cryptor.HmacSha256WithBase64(timestamp+"\n"+d.secret, d.secret)
		sign := url.QueryEscape(hmacCode)

		dingdingUrl := fmt.Sprintf("https://oapi.dingtalk.com/robot/send?access_token=%s&timestamp=%s&sign=%s", d.accessToken, timestamp, sign)
		postData := map[string]interface{}{
			"msgtype": "markdown",
			"markdown": map[string]string{
				"title": title,
				"text":  content,
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

		if respJson["errcode"] != float64(0) {
			logger.Log.Error(fmt.Sprintf("请求钉钉失败,errcode:%d,errmsg:%s", respJson["errcode"], respJson["errmsg"]))
			return
		}
	}()
}
