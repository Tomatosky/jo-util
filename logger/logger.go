package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
	"os"
)

var Log *zap.Logger

func InitLog(w io.Writer) {
	var coreArr []zapcore.Core

	//日志级别
	debugLevel := zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
		return lev <= zap.FatalLevel && lev >= zap.DebugLevel
	})
	infoLevel := zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
		return lev <= zap.FatalLevel && lev >= zap.InfoLevel
	})

	consoleEncoderConfig := zap.NewProductionEncoderConfig()
	consoleEncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05.000") //指定时间格式
	consoleEncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	consoleEncoder := zapcore.NewConsoleEncoder(consoleEncoderConfig)
	colorEncoder := &ColorEncoder{Encoder: consoleEncoder}
	consoleCore := zapcore.NewCore(colorEncoder, zapcore.AddSync(os.Stdout), debugLevel)
	coreArr = append(coreArr, consoleCore)

	if w != nil {
		fileEncoderConfig := zap.NewProductionEncoderConfig()
		fileEncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05.000")
		fileEncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
		fileEncoder := zapcore.NewJSONEncoder(fileEncoderConfig)
		fileWriteSyncer := zapcore.AddSync(w)
		fileCore := zapcore.NewCore(fileEncoder, fileWriteSyncer, infoLevel)
		coreArr = append(coreArr, fileCore)
	}

	Log = zap.New(zapcore.NewTee(coreArr...), zap.AddCaller(), zap.AddStacktrace(zap.ErrorLevel))
}
