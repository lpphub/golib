package zlog

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"time"
)

const (
	KeyCtxLogId    = "ctx_logId"
	KeyHeaderLogId = "X-Trace-logId"
)

func GetLogId(ctx *gin.Context) string {
	if ctx == nil {
		return generateLogId()
	}
	if logId := ctx.GetString(KeyCtxLogId); logId != "" {
		return logId
	}
	// 尝试从header中获取
	var logId string
	if ctx.Request != nil && ctx.Request.Header != nil {
		logId = ctx.GetHeader(KeyHeaderLogId)
	}
	if logId == "" {
		logId = generateLogId()
	}
	ctx.Set(KeyCtxLogId, logId)
	return logId
}

func SetHeaderLogId(ctx *gin.Context) {
	ctx.Header(KeyHeaderLogId, GetLogId(ctx))
}

func generateLogId() string {
	return strconv.FormatUint(uint64(time.Now().UnixNano())&0x7FFFFFFF|0x80000000, 10)
}
