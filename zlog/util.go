package zlog

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"time"
)

const (
	GinCtxLogId    = "_ctx_logId"
	GinModule      = "_log_module"
	GinHeaderLogId = "X-Trace-logId"
)

func GetLogId(ctx *gin.Context) string {
	if ctx == nil {
		return generateLogId()
	}
	if logId := ctx.GetString(GinCtxLogId); logId != "" {
		return logId
	}
	// 尝试从header中获取
	var logId string
	if ctx.Request != nil && ctx.Request.Header != nil {
		logId = ctx.GetHeader(GinHeaderLogId)
	}
	if logId == "" {
		logId = generateLogId()
	}
	ctx.Set(GinCtxLogId, logId)
	return logId
}

func GetModuleWithDefault(ctx *gin.Context, def string) string {
	if module := ctx.GetString(GinModule); module != "" {
		return module
	}
	ctx.Set(GinModule, def)
	return def
}

func SetHeaderLogId(ctx *gin.Context) {
	ctx.Header(GinHeaderLogId, GetLogId(ctx))
}

func generateLogId() string {
	return strconv.FormatUint(uint64(time.Now().UnixNano())&0x7FFFFFFF|0x80000000, 10)
}
