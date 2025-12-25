# jo-util

[![Go Version](https://img.shields.io/badge/Go-1.25-blue.svg)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

一个功能丰富的 Go 语言工具库,提供常用的工具函数和数据结构实现。

## 功能特性

### 数据结构 (mapUtil)
- **ConcurrentHashMap** - 线程安全的哈希映射
- **ConcurrentSkipListMap** - 基于跳表的并发有序映射
- **OrderedMap** - 保持插入顺序的映射
- **TreeMap** - 基于红黑树的有序映射
- **IMap** - 统一的 Map 接口,支持 JSON/BSON 序列化

### 集合工具 (setUtil)
- **HashSet** - 基于 map 实现的集合
- **ConcurrentHashSet** - 线程安全的集合

### 切片工具 (sliceUtil)
- **CopyOnWriteSlice** - 写时复制切片,适合读多写少场景
- 切片操作工具函数

### 队列工具 (queueUtil)
- **Queue** - 基础队列实现
- **SafeFixedQueue** - 线程安全的固定大小队列

### 缓存工具 (cacheUtil)
- 提供内存缓存功能
- 支持过期时间设置

### 加密工具 (cryptor)
- **AES 加密** - 支持 ECB、CBC、CTR、CFB、OFB 模式
- **DES 加密** - 支持 ECB、CBC、CTR、CFB、OFB 模式
- **RSA 加密** - 支持 PKCS1v15 和 OAEP 模式
- 多种填充方式 - NoPadding、ZeroPadding、Pkcs7Padding

### 日志工具 (logger)
- 基于 zap 的高性能日志库
- 彩色控制台输出
- 文件日志轮转 (基于 lumberjack)
- 日志监听器功能

### 线程池 (poolUtil)
- **AntsPool** - 基于 ants 的协程池封装
- **IdPool** - ID 池管理

### HTTP 工具 (httpUtil)
- HTTP 请求封装
- 简化的 API 调用

### 时间工具 (dateUtil)
- **TimeInterval** - 时间间隔计算
- 日期格式化和解析工具

### 字符串工具 (strUtil)
- 常用字符串处理函数

### 数字工具 (numberUtil)
- 数值处理和转换工具

### 随机工具 (randomUtil)
- 随机数生成
- 随机字符串生成

### ID 生成器 (idUtil)
- **UUID** - 基于 google/uuid
- **Snowflake** - 分布式 ID 生成 (基于 bwmarrin/snowflake)

### 文件工具 (fileUtil)
- 文件读写操作
- 文件路径处理

### 类型转换 (convertor)
- 各种类型之间的转换工具

### 系统工具 (osUtil)
- 操作系统相关的工具函数

### 事件工具 (eventUtil)
- 事件发布订阅机制

### MongoDB 工具 (mongoUtil)
- MongoDB 索引管理
- 数据库操作封装

## 安装

```bash
go get github.com/Tomatosky/jo-util
```

## 快速开始

### Map 使用示例

```go
package main

import (
    "fmt"
    "github.com/Tomatosky/jo-util/mapUtil"
)

func main() {
    // 创建并发安全的 HashMap
    m := mapUtil.NewConcurrentHashMap[string, int]()

    // 添加元素
    m.Put("key1", 100)
    m.Put("key2", 200)

    // 获取元素
    value := m.Get("key1")
    fmt.Println(value) // 输出: 100

    // 检查键是否存在
    exists := m.ContainsKey("key2")
    fmt.Println(exists) // 输出: true

    // 遍历
    m.Range(func(key string, value int) bool {
        fmt.Printf("%s: %d\n", key, value)
        return true // 返回 true 继续遍历
    })
}
```

### 加密工具示例

```go
package main

import (
    "fmt"
    "github.com/Tomatosky/jo-util/cryptor"
)

func main() {
    // AES 加密示例
    key := []byte("1234567890123456") // 16 字节密钥
    data := []byte("Hello, World!")

    // 加密
    encrypted := cryptor.AesEcbEncrypt(data, key, cryptor.Pkcs7Padding)
    fmt.Printf("加密后: %x\n", encrypted)

    // 解密
    decrypted := cryptor.AesEcbDecrypt(encrypted, key, cryptor.Pkcs7Padding)
    fmt.Printf("解密后: %s\n", decrypted)
}
```

### 日志工具示例

```go
package main

import (
    "github.com/Tomatosky/jo-util/logger"
)

func main() {
    // 简单初始化,日志会输出到文件和控制台
    log := logger.SimplyInit("./logs/app.log")

    log.Info("应用启动成功")
    log.Error("发生错误",
        zap.String("module", "main"),
        zap.Int("code", 500),
    )
}
```

### 线程池示例

```go
package main

import (
    "fmt"
    "github.com/Tomatosky/jo-util/poolUtil"
    "time"
)

func main() {
    // 创建大小为 10 的协程池
    pool := poolUtil.NewAntsPool(10)
    defer pool.Release()

    // 提交任务
    for i := 0; i < 100; i++ {
        taskID := i
        pool.Submit(func() {
            fmt.Printf("执行任务 %d\n", taskID)
            time.Sleep(100 * time.Millisecond)
        })
    }
}
```

### ID 生成器示例

```go
package main

import (
    "fmt"
    "github.com/Tomatosky/jo-util/idUtil"
)

func main() {
    // 生成 UUID
    uuid := idUtil.NewUUID()
    fmt.Println("UUID:", uuid)

    // 生成 Snowflake ID
    node, _ := idUtil.NewSnowflakeNode(1)
    id := node.Generate()
    fmt.Println("Snowflake ID:", id)
}
```

## 项目结构

```
jo-util/
├── cacheUtil/          # 缓存工具
├── convertor/          # 类型转换
├── cryptor/            # 加密工具
├── dateUtil/           # 时间工具
├── eventUtil/          # 事件工具
├── fileUtil/           # 文件工具
├── httpUtil/           # HTTP 工具
├── idUtil/             # ID 生成器
├── logger/             # 日志工具
├── mapUtil/            # Map 数据结构
├── mongoUtil/          # MongoDB 工具
├── numberUtil/         # 数字工具
├── osUtil/             # 系统工具
├── poolUtil/           # 线程池
├── queueUtil/          # 队列工具
├── randomUtil/         # 随机工具
├── setUtil/            # 集合工具
├── sliceUtil/          # 切片工具
└── strUtil/            # 字符串工具
```

## 依赖项

- [go.uber.org/zap](https://github.com/uber-go/zap) - 高性能日志库
- [gopkg.in/natefinch/lumberjack.v2](https://github.com/natefinch/lumberjack) - 日志轮转
- [github.com/panjf2000/ants/v2](https://github.com/panjf2000/ants) - 协程池
- [github.com/google/uuid](https://github.com/google/uuid) - UUID 生成
- [github.com/bwmarrin/snowflake](https://github.com/bwmarrin/snowflake) - Snowflake ID 生成
- [go.mongodb.org/mongo-driver](https://github.com/mongodb/mongo-go-driver) - MongoDB 驱动

## 测试

运行所有测试:

```bash
go test ./...
```

运行特定包的测试:

```bash
go test ./mapUtil
go test ./cryptor
```

## 贡献

欢迎提交 Pull Request 或 Issue!

## 许可证

本项目采用 MIT 许可证 - 详见 [LICENSE](LICENSE) 文件

## 作者

- GitHub: [@Tomatosky](https://github.com/Tomatosky)

## 更新日志

查看 [Releases](https://github.com/Tomatosky/jo-util/releases) 页面了解版本更新历史。
