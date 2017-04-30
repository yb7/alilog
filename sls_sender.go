package alilog

import (
  "time"
  "github.com/denverdino/aliyungo/sls"
  "os"
  "net"
  "crypto/md5"
  "encoding/hex"
  "github.com/golang/protobuf/proto"
)

var logChan = make(chan *logDto, 1000)

type logDto struct {
  Project  string
  LogStore string
  Time     time.Time
  Contents map[string]string
}

func (l *logDto) SlsLogContents() []*sls.Log_Content {
  c := make([]*sls.Log_Content, 0)
  for k, v := range l.Contents {
    c = append(c, &sls.Log_Content{
      Key:   proto.String(k),
      Value: proto.String(v),
    })
  }
  return c
}

func init() {
  go readLog()
}
func readLog() {
  flushTimer := time.NewTimer(time.Second * 3)
  ip := ipAddr()
  topic := ""
  const BUF_CAP = 100
  var buf = make([]*logDto, 0, BUF_CAP)
  for slsClient != nil {
    select {
    case <-flushTimer.C:
      _debug("time out of flush time, buf.size=%d", len(buf))
      if len(buf) > 0 {
        writeLogToSls(ip, topic, buf)
        buf = make([]*logDto, 0, BUF_CAP)
      }
    case msg := <-logChan:
      buf = append(buf, msg)
      if len(buf) >= BUF_CAP {
        writeLogToSls(ip, topic, buf)
        buf = make([]*logDto, 0, BUF_CAP)
      }
    }
  }
}

type logStoreKey struct {
  project string
  logStore string
}
func writeLogToSls(ip, topic string, buf []*logDto) {
  dividedByLogStore := make(map[logStoreKey][]*sls.Log)
  for _, dto := range buf {
    key := logStoreKey{dto.Project, dto.LogStore}
    logs, ok := dividedByLogStore[key]
    if !ok {
      logs = make([]*sls.Log, 0)
    }
    logs = append(logs, &sls.Log{
      Time:     proto.Uint32(uint32(dto.Time.Unix())),
      Contents: dto.SlsLogContents(),
    })
    dividedByLogStore[key] = logs
  }
  for key, value := range dividedByLogStore {
    writeLogToSlsStore(ip, topic, key.project, key.logStore, value)
  }

}
func writeLogToSlsStore(ip, topic, project, logStore string, logItems []*sls.Log) {
  logGroup := sls.LogGroup{
    Source: &ip,
    Topic:  &topic,
    Logs:   logItems,
  }

  go func() {
    _debug("write to sls >>>>> \n")
    req := &sls.PutLogsRequest{
      Project:  project,
      LogStore: logStore,
      LogItems: logGroup,
      HashKey:  getMD5Hash(ip),
    }
    for _, item := range req.LogItems.Logs {
      _debug("log at %s\n", *item.Time)
      for _, c := range item.Contents {
        _debug("    %s -> %s\n", *c.Key, *c.Value)
      }
    }
    err := slsClient.PutLogs(req)
    if err != nil {
      _debug("write to sls >>>>> %s\n", err.Error())
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
