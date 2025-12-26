package logger

import (
	"io"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var Log *zap.Logger

func init() {
	Log = InitLog(nil)
}

func InitLog(w io.Writer) *zap.Logger {
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

	return zap.New(zapcore.NewTee(coreArr...), zap.AddCaller(), zap.AddStacktrace(zap.ErrorLevel))
}

func SimplyInit(filename string) *zap.Logger {
	return InitLog(&lumberjack.Logger{
		Filename:   filename, //日志文件存放目录，如果文件夹不存在会自动创建
		MaxSize:    50,       //文件大小限制,单位MB
		MaxBackups: 10,       //最大保留日志文件数量
		MaxAge:     30,       //日志文件保留天数
		Compress:   false,    //是否压缩处理
		LocalTime:  true,
	})
}
