package alilog

import (
  "fmt"
  "os"
)

var defaultProjectName string
var defaultLogStore string

func init() {
  defaultProjectName = os.Getenv("ALILOG_PROJECT_NAME")
  defaultLogStore = os.Getenv("ALILOG_LOG_STORE")
  if len(defaultProjectName) == 0 || len(defaultLogStore) == 0{
    fmt.Printf("ALILOG_PROJECT_NAME / ALILOG_LOG_STORE shouldn't be empty, when use default logger")
  }
}

func defaultLogger() *SLog {
  return &SLog{
    projectName:  defaultProjectName,
    logStoreName: defaultLogStore,
    params: make(map[string]string),
  }
}
func LogWith(k, v string) *SLog {
	return defaultLogger().With(k, v)
}
func Tracef(format string, v ...interface{}) {
  defaultLogger().doLog("trace", format, v...)
}
func Debugf(format string, v ...interface{}) {
  defaultLogger().doLog("debug", format, v...)
}

func Infof(format string, v ...interface{}) {
  defaultLogger().doLog("info", format, v...)
}

func Warnf(format string, v ...interface{}) {
  defaultLogger().doLog("warn", format, v...)
}

func Errorf(format string, v ...interface{}) error {
  defaultLogger().doLog("error", format, v...)
	return fmt.Errorf(format, v...)
}
func Error(err error) error {
	if err != nil {
    defaultLogger().doLog("error", err.Error())
	}
	return err
}
func Fatal(err error) {
  defaultLogger().Error(err)
  os.Exit(1)
}
func Fatalf(format string, v ...interface{}) {
  defaultLogger().doLog("error", format, v...)
  os.Exit(1)
}
