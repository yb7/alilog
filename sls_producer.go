package alilog

import (
	"github.com/aliyun/aliyun-log-go-sdk/producer"
)

var producerInstance *producer.Producer

func StartSlsLog() {
	if !slsConfig.IsValid() {
		stdInfo.Println("sls config has not been inited")
		return
	}
	// https://github.com/aliyun/aliyun-log-go-sdk/blob/7cb5da33b5282aa0ea14bf56697ded4ce5818d37/producer/README.md?plain=1
	producerConfig := producer.GetDefaultProducerConfig()
	// 单个 producer 实例能缓存的日志大小上限，这里设置为5MB。
	producerConfig.TotalSizeLnBytes = 5 * 1024 * 1024
	// 单个 producer 实例能缓存的日志大小上限，这里设置为50k。
	producerConfig.MaxBatchSize = 50 * 1024
	// 当一个 ProducerBatch 中缓存的日志条数大于等于 batchCountThreshold 时，该 batch 将被发送，默认为 4096，这里设置成20条。
	producerConfig.MaxBatchCount = 20
	// 设置日志输出级别，默认值是Info,consumer中一共有4种日志输出级别，分别为debug,info,warn和error。
	producerConfig.AllowLogLevel = "debug"

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
