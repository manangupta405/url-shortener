package utils

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

func CustomLogFormatter(param gin.LogFormatterParams) string {
	var statusColor, methodColor, resetColor string
	if param.IsOutputColor() {
		statusColor = param.StatusCodeColor()
		methodColor = param.MethodColor()
		resetColor = param.ResetColor()
	}
	param.Latency = param.Latency / time.Microsecond
	level := "info"
	if param.StatusCode >= 400 {
		level = "error"
	} else if param.StatusCode >= 300 {
		level = "warn"
	}
	if param.ErrorMessage != "" {
		level = "error"
	}

	logData := map[string]interface{}{
		"type":      "GIN",
		"timestamp": param.TimeStamp.Format("2006/01/02 - 15:04:05"),
		"status":    fmt.Sprintf("%s%d%s", statusColor, param.StatusCode, resetColor),
		"latency":   param.Latency,
		"client_ip": param.ClientIP,
		"method":    fmt.Sprintf("%s%s%s", methodColor, param.Method, resetColor),
		"path":      param.Path,
		"error":     param.ErrorMessage,
		"level":     level,
	}

	jsonData, err := json.Marshal(logData)
	if err != nil {
		return fmt.Sprintf("Error marshaling log data: %v", err)
	}
	return string(jsonData) + "\n"
}
