package alilog

import (
  "os"
  "github.com/denverdino/aliyungo/common"
  "fmt"
  "github.com/denverdino/aliyungo/sls"
  "strings"
  "io/ioutil"
  "encoding/json"
)

var slsClient *sls.Client

type SlsConfig struct {
  AccessKeyID string
  AccessKeySecret string
  Region string
  EndPoint string
}


func readConfig(file string) SlsConfig {
  data, err := ioutil.ReadFile(file)
  if err != nil {
    panic(fmt.Errorf("error[%s] when read sls config file: %s", err.Error(), file))
  }
  var slsConfig SlsConfig
  err = json.Unmarshal(data, &slsConfig)
  if err != nil {
    panic(fmt.Errorf("error[%s] when unmarshal sls config\n %s", err.Error(), string(data)))
  }
  return slsConfig
}
func init() {
  cfgFile := os.Getenv("ALILOG_CONFIG")
  if len(cfgFile) == 0{
    stdInfo.Println("missing ALILOG_CONFIG sls start up failed")
    return
  }
  if _, err := os.Stat(cfgFile); os.IsNotExist(err) {
    panic(fmt.Sprintf("sls config file[%s], not exist", cfgFile))
  }
  slsConfig := readConfig(cfgFile)
  assertNotEmpty("slsConfig.AccessKeyID", slsConfig.AccessKeyID)
  assertNotEmpty("slsConfig.AccessKeySecret", slsConfig.AccessKeySecret)
  assertNotEmpty("slsConfig.Region", slsConfig.Region)
  slsRegion := assertRegion(slsConfig.Region)

  slsClient = sls.NewClientWithEndpoint(slsConfig.EndPoint, slsRegion, false,
    slsConfig.AccessKeyID, slsConfig.AccessKeySecret)

  stdInfo.Println("success create sls client")
}

func assertNotEmpty(key, value string) {
  if len(strings.TrimSpace(value)) == 0 {
    panic(fmt.Errorf("%s is empty", key))
  }
}

func assertRegion(slsRegion string) (reg common.Region) {
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
    panic(fmt.Errorf("not a valid aliyun region: [%s]", slsRegion))
  }
  return
}
