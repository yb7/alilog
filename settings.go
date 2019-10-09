package alilog

import (
	"os"
	"strings"
)

var ALI_INTERNAL_DEBUG = false

// debug 是 1，以后自增长，预留 0 给 trace
var _aliLogLevel = 1

func init() {
	if os.Getenv("ALI_INTERNAL_DEBUG") == "true" {
		stdInfo.Println("ALI_INTERNAL_DEBUG enabled")
		ALI_INTERNAL_DEBUG = true
	}
	_aliLogLevel = logLevelNum(os.Getenv("ALI_LOG_LEVEL"))
}

func ShouldLog(level string) bool {
	return logLevelNum(level) >= _aliLogLevel
}
func logLevelNum(level string) int {
	switch strings.ToUpper(level) {
  case "TRACE":
    return 0
  case "DEBUG":
    return 1
	case "INFO":
		return 2
	case "WARN":
		return 3
	case "ERROR":
		return 4
	default:
		return 1
	}
}
