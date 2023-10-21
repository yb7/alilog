```go
alilog.InitFromConfigFile(os.Getenv("ALILOG_CONFIG"))
alilog.StartSlsLog()
defer alilog.CloseSlsLog(3000)
```


## **producer配置详解**

[参考链接](https://github.com/aliyun/aliyun-log-go-sdk/blob/7cb5da33b5282aa0ea14bf56697ded4ce5818d37/producer/README.md?plain=1)

| 参数                | 类型   | 描述                                                         |
| ------------------- | ------ | ------------------------------------------------------------ |
| TotalSizeLnBytes    | Int64  | 单个 producer 实例能缓存的日志大小上限，默认为 100MB。       |
| MaxIoWorkerCount    | Int64  | 单个producer能并发的最多groutine的数量，默认为50，该参数用户可以根据自己实际服务器的性能去配置。 |
| MaxBlockSec         | Int    | 如果 producer 可用空间不足，调用者在 send 方法上的最大阻塞时间，默认为 60 秒。<br/>如果超过这个时间后所需空间仍无法得到满足，send 方法会抛出TimeoutException。如果将该值设为0，当所需空间无法得到满足时，send 方法会立即抛出 TimeoutException。如果您希望 send 方法一直阻塞直到所需空间得到满足，可将该值设为负数。 |
| MaxBatchSize        | Int64  | 当一个 ProducerBatch 中缓存的日志大小大于等于 batchSizeThresholdInBytes 时，该 batch 将被发送，默认为 512 KB，最大可设置成 5MB。 |
| MaxBatchCount       | Int    | 当一个 ProducerBatch 中缓存的日志条数大于等于 batchCountThreshold 时，该 batch 将被发送，默认为 4096，最大可设置成 40960。 |
| LingerMs            | Int64  | 一个 ProducerBatch 从创建到可发送的逗留时间，默认为 2 秒，最小可设置成 100 毫秒。 |
| Retries             | Int    | 如果某个 ProducerBatch 首次发送失败，能够对其重试的次数，默认为 10 次。<br/>如果 retries 小于等于 0，该 ProducerBatch 首次发送失败后将直接进入失败队列。 |
| MaxReservedAttempts | Int    | 每个 ProducerBatch 每次被尝试发送都对应着一个 Attemp，此参数用来控制返回给用户的 attempt 个数，默认只保留最近的 11 次 attempt 信息。<br/>该参数越大能让您追溯更多的信息，但同时也会消耗更多的内存。 |
| BaseRetryBackoffMs  | Int64  | 首次重试的退避时间，默认为 100 毫秒。 Producer 采样指数退避算法，第 N 次重试的计划等待时间为 baseRetryBackoffMs * 2^(N-1)。 |
| MaxRetryBackoffMs   | Int64  | 重试的最大退避时间，默认为 50 秒。                           |
| AdjustShargHash     | Bool   | 如果调用 send 方法时指定了 shardHash，该参数用于控制是否需要对其进行调整，默认为 true。 |
| Buckets             | Int    | 当且仅当 adjustShardHash 为 true 时，该参数才生效。此时，producer 会自动将 shardHash 重新分组，分组数量为 buckets。<br/>如果两条数据的 shardHash 不同，它们是无法合并到一起发送的，会降低 producer 吞吐量。将 shardHash 重新分组后，能让数据有更多地机会被批量发送。该参数的取值范围是 [1, 256]，且必须是 2 的整数次幂，默认为 64。 |
| AllowLogLevel       | String | 设置日志输出级别，默认值是Info,consumer中一共有4种日志输出级别，分别为debug,info,warn和error。 |
| LogFileName         | String | 日志文件输出路径，不设置的话默认输出到stdout。               |
| IsJsonType          | Bool   | 是否格式化文件输出格式，默认为false。                        |
| LogMaxSize          | Int    | 单个日志存储数量，默认为10M。                                |
| LogMaxBackups       | Int    | 日志轮转数量，默认为10。                                     |
| LogCompass          | Bool   | 是否使用gzip 压缩日志，默认为false。                         |
| Endpoint            | String | 服务入口，关于如何确定project对应的服务入口可参考文章[服务入口](https://help.aliyun.com/document_detail/29008.html?spm=a2c4e.11153940.blogcont682761.14.446e7720gs96LB)。 |
| AccessKeyID         | String | 账户的AK id。                                                |
| AccessKeySecret     | String | 账户的AK 密钥。                                              |
| NoRetryStatusCodeList  | []int  | 用户配置的不需要重试的错误码列表，当发送日志失败时返回的错误码在列表中，则不会重试。默认包含400，404两个值。                 |
| UpdateStsToken      | Func   | 函数类型，该函数内去实现自己的获取ststoken 的逻辑，producer 会自动刷新ststoken并放入client 当中。
| StsTokenShutDown    | channel| 关闭ststoken 自动刷新的通讯信道，当该信道关闭时，不再自动刷新ststoken值。当producer关闭的时候，该参数不为nil值，则会主动调用close去关闭该信道停止ststoken的自动刷新。 |
