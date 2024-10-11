package ware

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/lpphub/golib/zlog"
	"go.uber.org/zap"
	"io"
	"strings"
	"unsafe"
)

const _bodyLength = 1024

func LogTrace() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var (
			logId = zlog.GetLogId(ctx)

			reqBody  string
			respBody string
		)

		resp := &respWriter{body: bytes.NewBufferString(""), ResponseWriter: ctx.Writer}
		ctx.Writer = resp

		// 请求参数，涉及到回写，要在处理业务逻辑之前
		reqBody = getReqBody(ctx, _bodyLength)

		ctx.Next()

		if resp.body != nil {
			respBody = resp.body.String()
			if len(respBody) > _bodyLength {
				respBody = respBody[:_bodyLength]
			}
		}

		fields := []zap.Field{
			zap.String("logId", logId),
			zap.String("url", ctx.Request.URL.Path),
			zap.Int("status", resp.Status()),
			zap.String("request", reqBody),
			zap.String("response", respBody),
		}
		zlog.ZapLogger.Info("tracing...", fields...)
	}
}

type respWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w respWriter) WriteString(s string) (int, error) {
	if w.body != nil {
		w.body.WriteString(s)
	}
	return w.ResponseWriter.WriteString(s)
}

func (w respWriter) Write(b []byte) (int, error) {
	if w.body != nil {
		w.body.Write(b)
	}
	return w.ResponseWriter.Write(b)
}

// 请求参数
func getReqBody(c *gin.Context, maxReqBodyLen int) (reqBody string) {
	if maxReqBodyLen == -1 {
		return
	}

	if c.Request.Body != nil {
		if c.ContentType() == binding.MIMEMultipartPOSTForm {
			requestBody, err := c.GetRawData()
			if err != nil {
				zlog.Warn(c, "get http request body error: "+err.Error())
			}
			// 回写数据
			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
			if _, err := c.MultipartForm(); err != nil {
				zlog.Warn(c, "parse http request form body error: "+err.Error())
			}
			reqBody = c.Request.PostForm.Encode()
			// 回写数据
			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		} else if c.ContentType() == "application/octet-stream" {
			// ignore
		} else {
			requestBody, err := c.GetRawData()
			if err != nil {
				zlog.Warn(c, "get http request body error: "+err.Error())
			}
			reqBody = bytesToStr(requestBody)
			// 回写数据
			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}
	}

	if c.Request.URL.RawQuery != "" {
		reqBody += "&" + c.Request.URL.RawQuery
	}

	if len(reqBody) > maxReqBodyLen {
		reqBody = reqBody[:maxReqBodyLen]
	}
	return reqBody
}

func getCookie(ctx *gin.Context) string {
	cStr := ""
	for _, c := range ctx.Request.Cookies() {
		cStr += fmt.Sprintf("%s=%s&", c.Name, c.Value)
	}
	return strings.TrimRight(cStr, "&")
}

// converts byte slice to string without a memory allocation
func bytesToStr(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
