package alertUtil

import (
	"fmt"

	"github.com/Tomatosky/jo-util/httpUtil"
	"github.com/Tomatosky/jo-util/logger"
)

// GotifyAlert Gotify报警实现
type GotifyAlert struct {
	host  string
	token string
}

func (g *GotifyAlert) Alert(title string, content string) {
	go func() {
		postData := map[string]interface{}{
			"title":   title,
			"message": content,
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
