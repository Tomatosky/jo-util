package logger

import (
	"fmt"

	"go.uber.org/zap/buffer"
	"go.uber.org/zap/zapcore"
)

type ColorEncoder struct {
	zapcore.Encoder
}

func (e *ColorEncoder) EncodeEntry(ent zapcore.Entry, fields []zapcore.Field) (*buffer.Buffer, error) {
	buf, err := e.Encoder.EncodeEntry(ent, fields)
	if err != nil {
		return nil, err
	}

	color := getColor(ent.Level)
	reset := "\x1b[0m"

	// 在原始日志内容前后添加颜色代码
	colored := fmt.Sprintf("%s%s%s", color, buf.String(), reset)

	// 替换缓冲区内容
	newBuf := buffer.NewPool().Get()
	newBuf.AppendString(colored)
	buf.Free() // 释放原缓冲区

	return newBuf, nil
}

func getColor(level zapcore.Level) string {
	switch level {
	case zapcore.DebugLevel:
		return "\x1b[32m" // 绿色
	case zapcore.InfoLevel:
		return "\x1b[34m" // 蓝色
	case zapcore.WarnLevel:
		return "\x1b[33m" // 黄色
	case zapcore.ErrorLevel, zapcore.DPanicLevel, zapcore.PanicLevel, zapcore.FatalLevel:
		return "\x1b[31m" // 红色
	default:
		return "\x1b[0m"
	}
}
