package alilog

import (
	"fmt"
	"log"
	"strings"
	"time"

	sls "github.com/aliyun/aliyun-log-go-sdk"
)

type SLog struct {
	projectName  string
	logStoreName string
	project      *sls.LogProject
	logStore     *sls.LogStore
	params       map[string]string
}

func New(projectName, logStoreName string) *SLog {
	logProject := &sls.LogProject{
		Name:            projectName,
		Endpoint:        slsConfig.EndPoint,
		AccessKeyID:     slsConfig.AccessKeyID,
		AccessKeySecret: slsConfig.AccessKeySecret,
	}
	var retry_times int
	var logstore *sls.LogStore
	var err error
	var slsLog = &SLog{
		projectName:  projectName,
		logStoreName: logStoreName,
		project:      logProject,
	}
	fmt.Printf("[SLS] create log store to: Project[%s], LogStore[%s], Endpoint[%s]\n", projectName, logStoreName, slsConfig.EndPoint)
	for retry_times = 0; ; retry_times++ {
		if retry_times > 5 {
			return slsLog
		}
		logstore, err = logProject.GetLogStore(logStoreName)
		if err != nil {
			fmt.Printf("GetLogStore fail, retry:%d, err:%v\n", retry_times, err)
			if strings.Contains(err.Error(), sls.PROJECT_NOT_EXIST) {
				return slsLog
			} else if strings.Contains(err.Error(), sls.LOGSTORE_NOT_EXIST) {
				err = logProject.CreateLogStore(logStoreName, 1, 2)
				if err != nil {
					fmt.Printf("CreateLogStore fail, err: %s\n", err.Error())
				} else {
					fmt.Println("CreateLogStore success")
				}
			}
		} else {
			fmt.Printf("GetLogStore success, retry:%d, name: %s, ttl: %d, shardCount: %d, createTime: %d, lastModifyTime: %d\n", retry_times, logstore.Name, logstore.TTL, logstore.ShardCount, logstore.CreateTime, logstore.LastModifyTime)
			break
		}
		time.Sleep(200 * time.Millisecond)
	}
	slsLog.logStore = logstore
	return slsLog
	// return &SLog{
	// 	project:  project,
	// 	logStore: logStore,
	// }
}
func (l *SLog) With(k, v string) *SLog {
	params := make(map[string]string)
	for k, v := range l.params {
		params[k] = v
	}
	params[k] = v
	return &SLog{
		projectName:  l.projectName,
		logStoreName: l.logStoreName,
		project:      l.project,
		logStore:     l.logStore,
		params:       params,
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
	if ShouldLog(level) == false {
		return
	}
	msg := fmt.Sprintf(format, v...)
	params := make([]string, 0)

	fileParam := l.params["file"]
	if funcParam, ok := l.params["func"]; ok {
		fileParam = fmt.Sprintf("%s[%s]", fileParam, funcParam)
	}

	for k, v := range l.params {
		if k != "file" && k != "func" {
			params = append(params, fmt.Sprintf("%s[%s]", k, v))
		}
	}

	var stdLog *log.Logger
	switch level {
	case "debug":
		stdLog = stdDebug
	case "warn":
		stdLog = stdWarning
	case "error":
		stdLog = stdError
	default:
		stdLog = stdInfo
	}

	msgArr := make([]string, 0, 3)
	if len(fileParam) > 0 {
		msgArr = append(msgArr, fileParam)
	}
	if len(params) > 0 {
		msgArr = append(msgArr, strings.Join(params, ", "))
	}
	msgArr = append(msgArr, msg)
	stdLog.Println(strings.Join(msgArr, " - "))
	if l.logStore != nil {
		contents := map[string]string{
			"level":   level,
			"message": msg,
		}
		for k, v := range l.params {
			contents[k] = v
		}
		logChan <- &logDto{
			ProjectName:  l.projectName,
			LogStoreName: l.logStoreName,
			// Project:  l.project,
			LogStore: l.logStore,
			Time:     time.Now(),
			Contents: contents,
		}
	} else {
		_debug("logStore is null, ignore the log\n")
	}
}
