package web

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/lpphub/golib/logger/logx"
	"io"
	"strings"
	"time"
	"unsafe"
)

type AccessLogConfig struct {
	Enable    bool
	Module    string
	SkipPaths []string
}

const (
	_bodyLength = 2048
)

func AccessLog(conf AccessLogConfig) gin.HandlerFunc {
	var (
		skipPaths = conf.SkipPaths
		skipMap   = make(map[string]struct{})
	)
	if length := len(skipPaths); length > 0 {
		skipMap = make(map[string]struct{}, length)
		for _, path := range skipPaths {
			skipMap[path] = struct{}{}
		}
	}
	if conf.Module == "" {
		conf.Module = "app"
	}

	return func(ctx *gin.Context) {
		path := ctx.Request.URL.Path
		if _, ok := skipMap[path]; ok {
			return
		}

		start := time.Now()
		var (
			reqBody  string
			respBody string
		)
		resp := &respWriter{body: bytes.NewBufferString(""), ResponseWriter: ctx.Writer}
		ctx.Writer = resp

		// 请求参数，涉及到回写，要在处理业务逻辑之前
		reqBody = getReqBody(ctx, _bodyLength)

		ctx.Next()

		end := time.Now()

		if resp.body != nil {
			respBody = resp.body.String()
			if len(respBody) > _bodyLength {
				respBody = respBody[:_bodyLength]
			}
		}

		logx.FromGinCtx(ctx).Info().CallerSkipFrame(-1).
			Str("module", conf.Module).
			Str("url", path).
			Float64("cost_ms", getDiffTime(start, end)).
			Str("clientIp", getClientIp(ctx)).
			Int("status", resp.Status()).
			Str("request", reqBody).
			Str("response", respBody).
			Msg("access_log")
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
				logx.Err(c, err, "get http request body error")
			}
			// 回写数据
			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
			if _, err = c.MultipartForm(); err != nil {
				logx.Err(c, err, "parse http request form body error")
			}
			reqBody = c.Request.PostForm.Encode()
			// 回写数据
			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		} else if c.ContentType() == "application/octet-stream" {
			// ignore
		} else {
			requestBody, err := c.GetRawData()
			if err != nil {
				logx.Err(c, err, "get http request body error")
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

func getDiffTime(start, end time.Time) float64 {
	return float64(end.Sub(start).Nanoseconds()/1e4) / 100.0
}

func getClientIp(ctx *gin.Context) (clientIP string) {
	if ctx == nil {
		return clientIP
	}
	return ctx.ClientIP()
}

func SetHeaderLogId(ctx *gin.Context) {
	ctx.Header(logx.HeaderTraceId, logx.GetTraceId(ctx))
}
