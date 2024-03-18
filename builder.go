package alilog

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/pelletier/go-toml/v2"
)

var slsConfig *SlsConfig = nil

type SlsConfig struct {
	AccessKeyID     string
	AccessKeySecret string
	EndPoint        string
	ProjectName     string
	LogStore        string
	Topic           string
	Tags            map[string]string
}

func (s *SlsConfig) IsValid() bool {
	notValid := slsConfig == nil || len(slsConfig.AccessKeyID) == 0 || len(slsConfig.AccessKeySecret) == 0 || len(slsConfig.EndPoint) == 0
	return !notValid
}

func SetConfig(config *SlsConfig) {
	if config == nil {
		return
	}
	slsConfig = config

	if len(slsConfig.ProjectName) == 0 && len(os.Getenv("ALILOG_PROJECT_NAME")) > 0 {
		slsConfig.ProjectName = os.Getenv("ALILOG_PROJECT_NAME")
	}
	if len(slsConfig.LogStore) == 0 && len(os.Getenv("ALILOG_LOG_STORE")) > 0 {
		slsConfig.LogStore = os.Getenv("ALILOG_LOG_STORE")
	}
	if len(slsConfig.Tags) == 0 && len(os.Getenv("ALILOG_TAGS")) > 0 {
		tags := make(map[string]string)
		toml.Unmarshal([]byte(os.Getenv("ALILOG_TAGS")), tags)
		slsConfig.Tags = tags
	}
}

func InitFromConfigFile(cfgFile string) {
	if len(cfgFile) == 0 {

		stdInfo.Println("missing sls start up config file")
		return
	}
	stdInfo.Println("init througth ALILOG_CONFIG(config file)")
	if _, err := os.Stat(cfgFile); os.IsNotExist(err) {
		panic(fmt.Sprintf("sls config file[%s], not exist", cfgFile))
	}
	data, err := os.ReadFile(cfgFile)
	if err != nil {
		panic(fmt.Errorf("error[%s] when read sls config file: %s", err.Error(), cfgFile))
	}
	var slsConfigFromFile = &SlsConfig{}
	err = json.Unmarshal(data, slsConfigFromFile)
	if err != nil {
		panic(fmt.Errorf("error[%s] when unmarshal sls config\n %s", err.Error(), string(data)))
	}
	SetConfig(slsConfigFromFile)
	// stdInfo.Println("SLS CONFIG start")
	// stdInfo.Println(string(data))
	// stdInfo.Println("SLS CONFIG end")
}

func assertNotEmpty(key, value string) {
	if len(strings.TrimSpace(value)) == 0 {
		panic(fmt.Errorf("%s is empty", key))
	}
}
