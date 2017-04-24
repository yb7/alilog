package alilog

import (
  "os"
  "github.com/denverdino/aliyungo/common"
  "fmt"
  "github.com/denverdino/aliyungo/sls"
  "strings"
)

var slsClient *sls.Client
func init() {
  slsRegion, err := getSlsRegionFromEnv()
  if err != nil {
    return
  }

  accessKeyId := strings.TrimSpace(os.Getenv("SLS_ACCESSKEY_ID"))
  if len(accessKeyId) == 0 {
    return
  }
  accessKeySecret := strFromEnvNotEmpty("SLS_ACCESSKEY_SECRET")

  internalStr := os.Getenv("SLS_INTERNAL")
  var internal = false
  if internalStr == "true" {
    internal = true
  }

  slsClient = sls.NewClient(slsRegion, internal, accessKeyId, accessKeySecret)
  stdInfo.Println("success create sls client")
}

func getSlsRegionFromEnv() (reg common.Region, err error) {
  slsRegion := os.Getenv("SLS_REGION")
  switch slsRegion {
  case "cn-hangzhou": reg = common.Hangzhou
  case "cn-qingdao": reg = common.Qingdao
  case "cn-beijing": reg = common.Beijing
  case "cn-hongkong": reg = common.Hongkong
  case "cn-shenzhen": reg = common.Shenzhen
  case "cn-shanghai": reg = common.Shanghai
  case "cn-zhangjiakou": reg = common.Zhangjiakou

  case "ap-southeast-1": reg = common.APSouthEast1
  case "ap-northeast-1": reg = common.APNorthEast1
  case "ap-southeast-2": reg = common.APSouthEast2

  case "us-west-1": reg = common.USWest1
  case "us-east-1": reg = common.USEast1

  case "me-east-1": reg = common.MEEast1

  case "eu-central-1": reg = common.EUCentral1
  default:
    err = fmt.Errorf("not a valid aliyun region: [%s]", slsRegion)
  }
  return
}
func strFromEnvNotEmpty(key string) string {
  val := strings.TrimSpace(os.Getenv(key))
  if len(val) > 0 {
    return val
  }
  panic(fmt.Errorf("key[%s] in env, must not empty", key))
}
