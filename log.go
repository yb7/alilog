package alilog

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
	"time"
)

type SLog struct {
	projectName  string
	logStoreName string
	ip           string
	params       *LogParams
}

func New(projectName, logStoreName string) *SLog {
	logParams := &LogParams{
		TraceId: make([]string, 0),
		Params:  make(map[string]any),
	}
	slsLog := &SLog{
		projectName:  projectName,
		logStoreName: logStoreName,
		ip:           ipAddr(),
		params:       logParams,
	}
	return slsLog
}

type LogParams struct {
	TraceId []string
	Params  map[string]any
}

func (p *LogParams) getParams() map[string]any {
	if p.Params == nil {
		p.Params = make(map[string]any)
	}
	return p.Params
}

func (p *LogParams) getStrValue(k string) (string, bool) {
	v, ok := p.getParams()[k]
	if !ok {
		return "", false
	}
	str, _ := v.(string)
	return str, true
}
func (p *LogParams) appendTraceId(traceId ...string) {
	p.TraceId = slices.Compact(append(p.TraceId, traceId...))

}

func WithContext(ctx context.Context) *SLog {
	logParams, ok := getLogParamsFromContext(ctx)
	log := defaultLogger()
	if !ok {
		return log
	}
	log.params = logParams
	return log
}
func getLogParamsFromContext(ctx context.Context) (*LogParams, bool) {
	logParams, ok := ctx.Value("LOG_PARAMS").(*LogParams)
	return logParams, ok
}
func writeToContext(ctx context.Context, p *LogParams) context.Context {
	return context.WithValue(ctx, "LOG_PARAMS", p)
}
func (l *SLog) CreateLogContext(ctx context.Context) context.Context {
	return writeToContext(ctx, l.params)
}
func (l *SLog) WithTraceId(traceId ...string) *SLog {
	l.params.appendTraceId(traceId...)
	return l
}
func (l *SLog) With(k string, v any) *SLog {
	params := make(map[string]any)
	for k, v := range l.params.getParams() {
		params[k] = v
	}
	params[k] = v
	return &SLog{
		projectName:  l.projectName,
		logStoreName: l.logStoreName,
		ip:           l.ip,
		//project:      l.project,
		//logStore:     l.logStore,
		params: &LogParams{Params: params},
	}
}
func (l *SLog) Tracef(format string, v ...interface{}) {
	l.doLog("trace", format, v...)
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

func (l *SLog) Errorf(format string, v ...interface{}) error {
	l.doLog("error", format, v...)
	return fmt.Errorf(format, v...)
}
func (l *SLog) Error(err error) error {
	if err != nil {
		l.doLog("error", err.Error())
	}
	return err
}
func (l *SLog) Fatal(err error) {
	l.Error(err)
	os.Exit(1)
}
func (l *SLog) Fatalf(format string, v ...interface{}) {
	l.doLog("error", format, v...)
	os.Exit(1)
}

func (l *SLog) doLog(level string, format string, v ...interface{}) {
	if ShouldLog(level) == false {
		return
	}
	msg := format
	if len(v) > 0 {
		msg = fmt.Sprintf(format, v...)
	}

	params := make([]string, 0)

	fileName, fileNameExisted := l.params.getStrValue("file")

	funcName, funcNameExisted := l.params.getStrValue("func")

	var lineNumber = ""

	if !fileNameExisted || !funcNameExisted {
		// pc, file, line, ok := runtime.Caller(2)
		_, file, line, ok := runtime.Caller(2)
		if ok {
			if !fileNameExisted {
				arr := strings.Split(filepath.ToSlash(file), "/")
				fileName = arr[len(arr)-1]
				if len(arr) > 1 {
					fileName = fmt.Sprintf("%s/%s", arr[len(arr)-2], fileName)
				}
			}
			lineNumber = fmt.Sprintf(":%d", line)
		}
	}

	// fileParam := fmt.Sprintf("%s%s[%s]", fileName, lineNumber, funcName)

	fileParam := fmt.Sprintf("%s%s", fileName, lineNumber)

	for k, v := range l.params.getParams() {
		if k != "file" && k != "func" && k != "line" {
			params = append(params, fmt.Sprintf("%s[%s]", k, v))
		}
	}

	var stdLog *log.Logger
	switch level {
	case "trace":
		stdLog = stdTrace
	case "debug":
		stdLog = stdDebug
	case "info":
		stdLog = stdInfo
	case "warn":
		stdLog = stdWarning
	case "error":
		stdLog = stdError
	default:
		stdLog = stdDebug
	}

	msgArr := make([]string, 0, 3)
	msgArr = append(msgArr, fileParam)
	if len(params) > 0 {
		msgArr = append(msgArr, strings.Join(params, ", "))
	}
	msgArr = append(msgArr, msg)
	stdLog.Println(strings.Join(msgArr, " - "))
	if producerInstance != nil && len(l.projectName) > 0 && len(l.logStoreName) > 0 {

		contents := map[string]string{
			"level":   level,
			"message": msg,
		}

		topic := ""
		if slsConfig.Tags != nil {
			for k, v := range slsConfig.Tags {
				if k == "topic" {
					topic = v
				} else {
					contents[k] = v
				}
			}
		}

		if _, ok := contents["file"]; !ok {
			contents["file"] = fileName
		}
		if _, ok := contents["func"]; !ok {
			contents["func"] = funcName
		}
		if _, ok := contents["lineNumber"]; !ok {
			contents["lineNumber"] = lineNumber
		}
		if l.params != nil {
			if len(l.params.TraceId) > 0 {
				traceIdsBytes, _ := json.Marshal(l.params.TraceId)
				contents["traceId"] = string(traceIdsBytes)
			}

			for k, v := range l.params.getParams() {
				strValue, err := ToStringE(v)
				if err != nil {
					contents[k] = strValue
				}
			}
		}
		// fmt.Printf("print to %s, %s", l.projectName, l.logStoreName)
		ip := ipAddr()

		writeLogToSls(ip, topic, &logDto{
			ProjectName:  l.projectName,
			LogStoreName: l.logStoreName,
			Time:         time.Now(),
			Contents:     contents,
		})
		//logChan <- &logDto{
		//	ProjectName:  l.projectName,
		//	LogStoreName: l.logStoreName,
		//	Time:         time.Now(),
		//	Contents:     contents,
		//}
	} else {
		_debug("logStore is null, ignore the log\n")
	}
}
