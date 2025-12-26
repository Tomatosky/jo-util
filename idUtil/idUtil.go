package idUtil

import (
	"fmt"
	"strings"

	"github.com/Tomatosky/jo-util/logger"
	"github.com/bwmarrin/snowflake"
	"github.com/google/uuid"
)

// 全局snowflake节点（默认节点ID为1）
var snowflakeNode *snowflake.Node

// RandomUUID 生成标准格式UUID（带连字符）
func RandomUUID() string {
	return uuid.NewString()
}

// SimpleUUID 生成不带连字符的UUID
func SimpleUUID() string {
	return strings.ReplaceAll(uuid.NewString(), "-", "")
}

func getSnowflake() snowflake.ID {
	if snowflakeNode == nil {
		var err error
		snowflakeNode, err = snowflake.NewNode(1)
		if err != nil {
			logger.Log.Fatal(fmt.Sprintf("Snowflake node initialization failed: %v", err))
		}
	}
	return snowflakeNode.Generate()
}

// GetSnowflakeNextId 生成雪花算法ID
func GetSnowflakeNextId() int64 {
	return getSnowflake().Int64()
}

// GetSnowflakeNextIdStr 生成字符串格式的雪花算法ID
func GetSnowflakeNextIdStr() string {
	return getSnowflake().String()
}
