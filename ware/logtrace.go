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
	return func(c *gin.Context) {
		var (
			logId = zlog.GetLogId(c)

			reqBody  string
			respBody string
		)

		resp := &respWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = resp

		// 请求参数，涉及到回写，要在处理业务逻辑之前
		reqBody = getReqBody(c, _bodyLength)

		c.Next()

		if resp.body != nil {
			respBody = resp.body.String()
			if len(respBody) > _bodyLength {
				respBody = respBody[:_bodyLength]
			}
		}

		fields := []zap.Field{
			zap.String("logId", logId),
			zap.String("url", c.Request.URL.Path),
			zap.String("reqBody", reqBody),
			zap.String("respBody", respBody),
		}
		zlog.ZapLogger.Info("trace...", fields...)
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
		return reqBody
	}

	if c.Request.Body != nil && c.ContentType() == binding.MIMEMultipartPOSTForm {
		requestBody, err := c.GetRawData()
		if err != nil {
			zlog.Warn(c, "get http request body error: "+err.Error())
		}
		c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))

		if _, err := c.MultipartForm(); err != nil {
			zlog.Warn(c, "parse http request form body error: "+err.Error())
		}
		reqBody = c.Request.PostForm.Encode()

		// 回写参数
		c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))

	} else if c.Request.Body != nil && c.ContentType() == "application/octet-stream" {

	} else if c.Request.Body != nil {
		requestBody, err := c.GetRawData()
		if err != nil {
			zlog.Warn(c, "get http request body error: "+err.Error())
		}
		reqBody = string(requestBody)
		c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
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
