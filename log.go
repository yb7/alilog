package alilog

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

type SLog struct {
	projectName  string
	logStoreName string
	ip           string
	//project      *sls.LogProject
	//logStore     *sls.LogStore
	params map[string]string
}

//func getLogStore(projectName, logStoreName string) *sls.LogStore {
//	logKey := logStoreKey(projectName, logStoreName)
//	data, ok := logCache.Load(logKey)
//	if !ok {
//		return nil
//	}
//	return data.(*slsLogData).logStore
//}

//type slsLogData struct {
//	//once         *sync.Once
//	project  *sls.LogProject
//	logStore *sls.LogStore
//}
//
//func logStoreKey(projectName, logStoreName string) string {
//	return fmt.Sprintf("%s-%s", projectName, logStoreName)
//}

//var logCache sync.Map

//var mutex = &sync.Mutex{}

//func initSlsLogData(projectName, logStoreName string) {
//	logKey := logStoreKey(projectName, logStoreName)
//	//const MaxRetry = 5
//	doOnce(logKey, func() error {
//		logProject := &sls.LogProject{
//			Name:            projectName,
//			Endpoint:        slsConfig.EndPoint,
//			AccessKeyID:     slsConfig.AccessKeyID,
//			AccessKeySecret: slsConfig.AccessKeySecret,
//		}
//
//		fmt.Printf("[SLS] create log store to: Project[%s], LogStore[%s], Endpoint[%s]\n", projectName, logStoreName, slsConfig.EndPoint)
//
//		//var retry_times int
//		for retry_times := 0; retry_times < 5 && len(slsConfig.EndPoint) > 0; retry_times++ {
//			logstore, err := logProject.GetLogStore(logStoreName)
//			if err != nil {
//				fmt.Printf("GetLogStore fail, retry:%d, err:%v\n", retry_times, err)
//				if strings.Contains(err.Error(), sls.PROJECT_NOT_EXIST) {
//					return err
//				} else if strings.Contains(err.Error(), sls.LOGSTORE_NOT_EXIST) {
//					err = logProject.CreateLogStore(logStoreName, 30, 2, true, 100)
//					if err != nil {
//						fmt.Printf("CreateLogStore fail, err: %s\n", err.Error())
//					} else {
//						fmt.Println("CreateLogStore success")
//					}
//				}
//			} else {
//				fmt.Printf("GetLogStore success, retry:%d, name: %s, ttl: %d, shardCount: %d, createTime: %d, lastModifyTime: %d\n",
//					retry_times, logstore.Name, logstore.TTL, logstore.ShardCount, logstore.CreateTime, logstore.LastModifyTime)
//				logCache.Store(logKey, &slsLogData{
//					logProject, logstore,
//				})
//				return nil
//			}
//			time.Sleep(200 * time.Millisecond)
//		}
//		return errors.New("[SLS] exceed max retry, create log store failed")
//		//slsLog.logStore = logstore
//	})
//}

func New(projectName, logStoreName string) *SLog {
	slsLog := &SLog{
		projectName:  projectName,
		logStoreName: logStoreName,
		ip:           ipAddr(),
		params:       make(map[string]string),
	}
	//go initSlsLogData(projectName, logStoreName)

	return slsLog
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
		ip:           l.ip,
		//project:      l.project,
		//logStore:     l.logStore,
		params: params,
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

	msg := fmt.Sprintf(format, v...)
	params := make([]string, 0)

	fileName, fileNameExisted := l.params["file"]
	funcName, funcNameExisted := l.params["func"]

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
			// if !funcNameExisted {
			// 	funcName = runtime.FuncForPC(pc).Name()
			// 	firstSlash := strings.LastIndex(funcName, "/")
			// 	if firstSlash > -1 {
			// 		funcName = funcName[firstSlash+1:]
			// 	}
			// 	if strings.Index(funcName, ".") > -1 {
			// 		fileNameFirstPart := strings.Split(fileName, "/")[0]
			// 		funcNameFirstPart := funcName[0:strings.Index(funcName, ".")]
			// 		if fileNameFirstPart == funcNameFirstPart {
			// 			funcName = funcName[strings.Index(funcName, ".")+1:]
			// 		}
			// 	}
			// }
			lineNumber = fmt.Sprintf(":%d", line)
		}
	}

	// fileParam := fmt.Sprintf("%s%s[%s]", fileName, lineNumber, funcName)

	fileParam := fmt.Sprintf("%s%s", fileName, lineNumber)

	for k, v := range l.params {
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
	if len(strings.TrimSpace(os.Getenv("ALILOG_CONFIG"))) > 0 && len(l.projectName) > 0 && len(l.logStoreName) > 0 {
		contents := map[string]string{
			"level":   level,
			"message": msg,
		}
		for k, v := range l.params {
			contents[k] = v
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
		ip := ipAddr()
		topic := ""
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
