package alilog

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	sls "github.com/aliyun/aliyun-log-go-sdk"
	"github.com/golang/protobuf/proto"
	"net"
	"os"
	"time"
)

type logDto struct {
	ProjectName  string
	LogStoreName string
	//LogStore     *sls.LogStore
	Time     time.Time
	Contents map[string]string
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

func writeLogToSls(ip, topic string, dto *logDto) {
	if producerInstance == nil {
		return
	}
	//dividedByLogStore := make(map[*sls.LogStore][]*sls.Log)
	err := producerInstance.SendLog(dto.ProjectName, dto.LogStoreName, topic, ip, &sls.Log{
		Time:     proto.Uint32(uint32(dto.Time.Unix())),
		Contents: dto.SlsLogContents(),
	})
	if err != nil {
		fmt.Println(err)
	}
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
