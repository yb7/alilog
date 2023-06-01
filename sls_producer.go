package alilog

import (
	"github.com/aliyun/aliyun-log-go-sdk/producer"
	"os"
	"os/signal"
)

var producerInstance *producer.Producer

func StartSlsLog() {
	if len(slsConfig.AccessKeyID) == 0 {
		stdInfo.Println("sls config has not been inited")
		return
	}
	producerConfig := producer.GetDefaultProducerConfig()
	producerConfig.Endpoint = slsConfig.EndPoint
	producerConfig.AccessKeyID = slsConfig.AccessKeyID
	producerConfig.AccessKeySecret = slsConfig.AccessKeySecret
	producerInstance = producer.InitProducer(producerConfig)
	ch := make(chan os.Signal)
	signal.Notify(ch, os.Kill, os.Interrupt)
	producerInstance.Start()
}

func SafeCloseSlsLog() {
	if producerInstance != nil {
		producerInstance.SafeClose()
	}
}
func CloseSlsLog(timeoutMs int64) {
	if producerInstance != nil {
		producerInstance.Close(timeoutMs)
	}
}
