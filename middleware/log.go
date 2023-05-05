package middleware

import (
	"fmt"
	"io"
	"os"

	"github.com/gin-gonic/gin"
)

func LoggingMiddleware(logFilePath string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		logger := gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
			logMsg := fmt.Sprintf("%s - %s [%s] \"%s %s %s\" %d %d \"%s\" \"%s\" %s\n",
				param.ClientIP,
				param.Request.Header.Get("X-Forwarded-For"),
				param.TimeStamp.Format("02/Jan/2006:15:04:05 -0700"),
				param.Method,
				param.Path,
				param.Request.Proto,
				param.StatusCode,
				param.BodySize,
				param.Request.Header.Get("Referer"),
				param.Request.Header.Get("User-Agent"),
				param.Latency,
			)

			file, err := os.OpenFile(logFilePath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
			if err != nil {
				panic(err)
			}
			defer file.Close()

			fileInfo, err := file.Stat()
			if err != nil {
				panic(err)
			}
			if fileInfo.Size() > 0 {
				_, err = file.Seek(0, io.SeekEnd)
				if err != nil {
					panic(err)
				}
			}

			if _, err := file.WriteString(logMsg); err != nil {
				panic(err)
			}

			return logMsg
		})
		logger(ctx)

		ctx.Next()
	}
}
