package main

import (
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
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

// 替换gin.Default中的Logger和Recover中间件
func GinLogger(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		c.Next()

		cost := time.Since(start)
		logger.Info(path,
			zap.Int("status", c.Writer.Status()),
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("query", query),
			zap.String("ip", c.ClientIP()),
			zap.String("user-agent", c.Request.UserAgent()),
			zap.String("errors", c.Errors.ByType(gin.ErrorTypePrivate).String()),
			zap.Duration("cost", cost),
		)
	}
}

// GinRecover recover项目可能出现的panic
func GinRecover(logger *zap.Logger, stack bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}
				httpRequest, _ := httputil.DumpRequest(c.Request, false)
				if brokenPipe {
					logger.Error(c.Request.URL.Path,
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
					// If the connection is dead, we can't write a status to it.
					c.Error(err.(error)) // nolint: errcheck
					c.Abort()
					return
				}

				if stack {
					logger.Error("[Recovery from panic]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
						zap.String("stack", string(debug.Stack())),
					)
				} else {
					logger.Error("[Recovery from panic]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
				}
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}

func main() {
	r := gin.New()
	r.Use(GinLogger(logger), GinRecover(logger, true))
	r.GET("/hello", func(c *gin.Context) {
		c.String(200, "hello gin")
	})
	r.Run()
}
