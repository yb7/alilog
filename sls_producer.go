package alilog

import (
	"github.com/aliyun/aliyun-log-go-sdk/producer"
)

var producerInstance *producer.Producer

func StartSlsLog() {
	if slsConfig == nil || len(slsConfig.AccessKeyID) == 0 {
		stdInfo.Println("sls config has not been inited")
		return
	}
	producerConfig := producer.GetDefaultProducerConfig()
	producerConfig.Endpoint = slsConfig.EndPoint
	producerConfig.AccessKeyID = slsConfig.AccessKeyID
	producerConfig.AccessKeySecret = slsConfig.AccessKeySecret
	producerInstance = producer.InitProducer(producerConfig)
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
