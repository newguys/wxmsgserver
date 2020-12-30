package log

import (
	"fmt"
	"os"
	"path"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

//Logger 日志
func Logger() *logrus.Logger {
	now := time.Now()
	logFilePath := ""
	if dir, err := os.Getwd(); err == nil {
		logFilePath = dir + "/logs/"
	}
	if err := os.MkdirAll(logFilePath, 0777); err != nil {
		fmt.Printf("log mkdir failed err:%s,logfilepath:%s", err.Error(), logFilePath)
	}
	logFileName := fmt.Sprintf("wxmsgserver%s.log", now.Format("2006-01-02"))

	//日志文件
	fileName := path.Join(logFilePath, logFileName)
	if _, err := os.Stat(fileName); err != nil {
		if _, err := os.Create(fileName); err != nil {
			fmt.Printf("log create file failed :%s", err.Error())
		}
	}

	//写入文件
	src, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		fmt.Printf("log openfile failed err :%s", err.Error())
	}

	logger := logrus.New()

	logger.Out = src
	logger.WriterLevel(logrus.DebugLevel)
	logger.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})
	return logger
}

//LoggerToFile 设置gin 日志middleware
func LoggerToFile() gin.HandlerFunc {
	logger := Logger()
	return func(c *gin.Context) {
		startTime := time.Now()
		c.Next()
		endTime := time.Now()

		latencyTime := endTime.Sub(startTime)
		reqMethod := c.Request.Method
		reqUrl := c.Request.RequestURI
		statusCode := c.Writer.Status()
		clientIp := c.ClientIP()
		logger.Infof("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
			clientIp,
			startTime,
			reqMethod,
			reqUrl,
			statusCode,
			latencyTime,
			c.Request.UserAgent(),
		)
	}
}
