package alilog

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net"
	"os"
	"time"
  "strconv"
	sls "github.com/aliyun/aliyun-log-go-sdk"
	"github.com/golang/protobuf/proto"
)

//var TOTAL_BUF_SIZE = 10000
const LOG_SENDER_TIMER = time.Second * 1

var bufCap = 100

var logChan chan *logDto

type logDto struct {
	ProjectName  string
	LogStoreName string
	//LogStore     *sls.LogStore
	Time         time.Time
	Contents     map[string]string
}

func (l *logDto) SlsLogContents() []*sls.LogContent {
	c := make([]*sls.LogContent, 0)
	for k, v := range l.Contents {
		c = append(c, &sls.LogContent{
			Key:   proto.String(k),
			Value: proto.String(v),
		})
	}
	return c
}

func init() {
  var bufCapStr = os.Getenv("ALILOG_BUF_CAP")
  if _bufCap, err := strconv.Atoi(bufCapStr); err != nil && _bufCap > 0 {
    bufCap = _bufCap
  }
  chanCap := bufCap*200
  if _chanCap, err := strconv.Atoi(os.Getenv("ALILOG_LOGCHAN_CAP")); err != nil && _chanCap > bufCap {
    chanCap = _chanCap
  }
  logChan = make(chan *logDto, chanCap)
	go readLog()
}
func readLog() {
	flushTimer := time.NewTicker(LOG_SENDER_TIMER)
	ip := ipAddr()
	topic := ""
	//const BUF_CAP = TOTAL_BUF_SIZE / 10
	var buf = make([]*logDto, 0, bufCap)
	for {
		select {
		case <-flushTimer.C:
			_debug("time out of flush time, buf.size=%d\n", len(buf))
			if len(buf) > 0 {
				writeLogToSls(ip, topic, buf)
				buf = make([]*logDto, 0, bufCap)
			}
		case msg := <-logChan:
			buf = append(buf, msg)
			if len(buf) >= bufCap {
				writeLogToSls(ip, topic, buf)
				buf = make([]*logDto, 0, bufCap)
			}
		}
	}
}

// type logStoreKey struct {
// 	project  string
// 	logStore string
// }

func writeLogToSls(ip, topic string, buf []*logDto) {
	dividedByLogStore := make(map[*sls.LogStore][]*sls.Log)
	for _, dto := range buf {
		// key := logStoreKey{dto.ProjectName, dto.LogStoreName}
		// dto.LogStore.
		logStore := getLogStore(dto.ProjectName, dto.LogStoreName)
		if logStore == nil {
		  continue
    }
		logs, ok := dividedByLogStore[logStore]
		if !ok {
			logs = make([]*sls.Log, 0)
		}
		logs = append(logs, &sls.Log{
			Time:     proto.Uint32(uint32(dto.Time.Unix())),
			Contents: dto.SlsLogContents(),
		})
		dividedByLogStore[logStore] = logs
	}
	_debug("divide logs to %s log stores\n", len(dividedByLogStore))
	for logStore, value := range dividedByLogStore {
		writeLogToSlsStore(ip, topic, logStore, value)
	}

}
func writeLogToSlsStore(ip, topic string, logStore *sls.LogStore, logItems []*sls.Log) {
	logGroup := &sls.LogGroup{
		Source: &ip,
		Topic:  &topic,
		Logs:   logItems,
	}

	go func() {
		_debug("write to sls >>>>> \n")
		err := logStore.PutLogs(logGroup)
		// req := &sls.PutLogsRequest{
		// 	Project:  project,
		// 	LogStore: logStore,
		// 	LogItems: logGroup,
		// 	HashKey:  getMD5Hash(ip),
		// }
		// for _, item := range req.LogItems.Logs {
		// 	_debug("log at %s\n", *item.Time)
		// 	for _, c := range item.Contents {
		// 		_debug("    %s -> %s\n", *c.Key, *c.Value)
		// 	}
		// }
		// err := slsClient.PutLogs(req)
		if err != nil {
			fmt.Printf("[SLS] error, write to sls >>>>> %s\n", err.Error())
		}
	}()
}

func ipAddr() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		os.Stderr.WriteString("Oops: " + err.Error() + "\n")
		return "unknown"
	}

	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return "unknown"
}
func getMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}
