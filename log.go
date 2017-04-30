package alilog

import (
  "time"
  "fmt"
  "strings"
  "log"
)

type SLog struct {
  project string
  logStore string
  params map[string]string
}

func New(project, logStore string) *SLog {
  return &SLog{
    project: project,
    logStore: logStore,
  }
}
func (l *SLog) With(k, v string) *SLog {
  params := make(map[string]string)
  for k, v := range l.params {
    params[k] = v
  }
  params[k] = v
  return &SLog{
    project: l.project,
    logStore: l.logStore,
    params: params,
  }
}
func (l *SLog) Debugf(format string, v ...interface{}) {
  l.doLog("debug", format, v...)
}

func (l *SLog) Infof(format string, v ...interface{}) {
  l.doLog("info", format, v...)
}

func (l *SLog) Warnf(format string, v ...interface{}) {
  l.doLog("warn", format, v...)
}

func (l *SLog) Errorf(format string, v ...interface{}) {
  l.doLog("error", format, v...)
}
func (l *SLog) Error(err error) {
  if err != nil {
    l.doLog("error", err.Error())
  }
}

func (l *SLog) doLog(level string, format string, v ...interface{}) {
  msg := fmt.Sprintf(format, v...)
  params := make([]string, 0)
  for k, v := range l.params {
    params = append(params, fmt.Sprintf("%s[%s]", k, v))
  }

  var stdLog *log.Logger
  switch level {
  case "debug": stdLog = stdDebug
  case "warn": stdLog = stdWarning
  case "error": stdLog = stdError
  default:
    stdLog = stdInfo
  }
  stdLog.Println(msg + " " + strings.Join(params, ", "))
  if slsClient != nil {
    contents := map[string]string {
      "level": level,
      "message": msg,
    }
    for k, v := range l.params {
      contents[k] = v
    }
    logChan <- &logDto{
      Project: l.project,
      LogStore: l.logStore,
      Time: time.Now(),
      Contents: contents,
    }
  }
}
