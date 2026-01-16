package alertUtil

import (
	"fmt"
	"strings"

	"github.com/Tomatosky/jo-util/logger"
)

// Alert 报警接口，由用户自行实现
type Alert interface {
	Alert(title string, content string)
}

// DefaultAlert 默认报警实现
type DefaultAlert struct{}

func (d *DefaultAlert) Alert(title string, content string) {
	content = strings.ReplaceAll(content, "\n", "")
	logger.Log.Warn(fmt.Sprintf("%s: %s", title, content))
}
