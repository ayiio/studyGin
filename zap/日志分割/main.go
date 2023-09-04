package main

import (
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"net/http"
)

var logger *zap.Logger
var sugarLogger *zap.SugaredLogger

func initLogger() {
	writeSyncer := getLogWriter()
	encoder := getEncoder()
	//zapcore.NewCore
	core := zapcore.NewCore(encoder, writeSyncer, zapcore.DebugLevel)

	//logger = zap.New(core)
	logger = zap.New(core, zap.AddCaller()) //添加函数调用信息
	sugarLogger = logger.Sugar()
}

// json格式encoder
func getEncoder() zapcore.Encoder {
	//zapcore.NewConsoleEncoder()  console打印格式
	//zap.NewProductionEncoderConfig()
	encodeConfig := zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder, //可读时间戳
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
	return zapcore.NewJSONEncoder(encodeConfig)
}

// writeWrapper
//func getLogWriter() zapcore.WriteSyncer {
//	//Open每次追加os.O_CREATE|os.O_RDWR|os.O_APPEND, 774， Create每次新创建
//	file, _ := os.Create("./test.log") //写出到test.log
//	return zapcore.AddSync(file)
//}

// zap本身不支持日志分割，借助第三方库实现 go get -u github.com/natefinch/lumberjack
func getLogWriter() zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   "./test.log",
		MaxSize:    1,     //M
		MaxBackups: 5,     //备份几份
		MaxAge:     30,    //day
		Compress:   false, //是否压缩
		LocalTime:  true,
	}
	return zapcore.AddSync(lumberJackLogger)
}

func simpleHttpGet(url string) {
	// logger
	resp, err := http.Get(url)
	if err != nil {
		logger.Error(
			"error fetching url...",
			zap.String("url", url),
			zap.Error(err))
	} else {
		logger.Info("success..",
			zap.String("statusCode", resp.Status),
			zap.String("url", url))
		resp.Body.Close()
	}

	// sugarlogger
	sugarLogger.Debugf("Trying to hit GET request for %s", url)
	resp, err = http.Get(url)
	if err != nil {
		sugarLogger.Errorf("Error fetching URL is : Error = %s", url, err)
	} else {
		sugarLogger.Infof("Success! statusCode = %s for URL %s", resp.Status, url)
		resp.Body.Close()
	}
}

func main() {
	initLogger()
	defer logger.Sync()

	for i := 0; i < 10000; i++ {
		logger.Info("test ....")
	}
	simpleHttpGet("www.baidu.com")
	simpleHttpGet("http://www.baidu.com")
}
