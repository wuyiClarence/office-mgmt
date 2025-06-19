package middleware

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"platform-backend/utils/log"
	"time"

	"github.com/gin-gonic/gin"
)

type responseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (r *responseBodyWriter) Write(b []byte) (int, error) {
	r.body.Write(b)                  // 写入缓冲区以记录响应体
	return r.ResponseWriter.Write(b) // 写回客户端
}

func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// 捕获响应体
		respBody := &bytes.Buffer{}
		c.Writer = &responseBodyWriter{body: respBody, ResponseWriter: c.Writer}

		// 读取请求体
		var reqBody []byte
		contentType := c.GetHeader("Content-Type")
		if contentType != "multipart/form-data" && c.Request.ContentLength <= 1<<20 { // 跳过文件类型和大于1MB的请求体
			if c.Request.Body != nil {
				reqBody, _ = ioutil.ReadAll(io.LimitReader(c.Request.Body, 1<<20)) // 限制读取1MB
				c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(reqBody))        // 重置Body供后续使用
			}
		}

		c.Next() // 继续处理请求

		// 日志信息
		end := time.Now()
		latency := end.Sub(start)
		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()
		if raw != "" {
			path = path + "?" + raw
		}

		// 日志输出
		_, _ = fmt.Fprintf(log.APILogger, "[GIN] %v | %3d | %13v | %15s | %-7s %s \n Request:[%s]\n Response:[%s]\n",
			end.Format("2006/01/02 - 15:04:05"),
			statusCode,
			latency,
			clientIP,
			method,
			path,
			string(reqBody),
			respBody.String(),
		)

		if gin.IsDebugging() {
			_, _ = fmt.Fprintf(os.Stdout, "[GIN] %v | %3d | %13v | %15s | %-7s %s\n",
				end.Format("2006/01/02 - 15:04:05"),
				statusCode,
				latency,
				clientIP,
				method,
				path,
			)
		}
	}
}
