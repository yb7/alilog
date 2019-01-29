package main

import (
	"encoding/json"
	"net/url"
	"time"

	"github.com/yb7/alilog"
)

func main() {
	log := alilog.New("dpns", "controller-push-dev")
	str := "AccessKeyId=LTAIsXIVTn5RFE9Y&Action=Push&AppKey=24650371&Body=%E3%80%90%E6%B5%B7%E6%A0%BC%E9%80%9A%E4%BF%A1%EF%BC%9A%E6%8E%A7%E8%82%A1%E8%82%A1%E4%B8%9C%E6%8B%9F%E5%A2%9E%E6%8C%81%E4%B8%8D%E8%B6%8510%E4%BA%BF%E5%85%83%E3%80%91%E6%B5%B7%E6%A0%BC%E9%80%9A%E4%BF%A1%E5%85%AC%E5%91%8A%EF%BC%8C%E5%85%AC%E5%8F%B8%E6%8E%A7%E8%82%A1%E8%82%A1%E4%B8%9C%E6%97%A0%E7%BA%BF%E7%94%B5%E9%9B%86%E5%9B%A2%E8%AE%A1%E5%88%92%E6%9C%AA%E6%9D%A56%E4%B8%AA%E6%9C%88%E5%86%85%E5%A2%9E%E6%8C%81%E5%85%AC%E5%8F%B8%E8%82%A1%E4%BB%BD%EF%BC%8C%E5%A2%9E%E6%8C%81%E9%87%91%E9%A2%9D%E6%8B%9F%E4%B8%8D%E8%B6%85%E8%BF%8710%E4%BA%BF%E5%85%83%EF%BC%8C%E5%A2%9E%E6%8C%81%E6%AF%94%E4%BE%8B%E4%B8%8D%E8%B6%85%E8%BF%87%E5%85%AC%E5%8F%B8%E6%80%BB%E8%82%A1%E6%9C%AC%E7%9A%845%25%EF%BC%8C%E5%A2%9E%E6%8C%81%E4%BB%B7%E6%A0%BC%E4%B8%8D%E9%AB%98%E4%BA%8E16%E5%85%83%2F%E8%82%A1%E3%80%82%E5%85%AC%E5%8F%B8%E6%9C%80%E6%96%B0%E8%82%A1...&DeviceType=ALL&Format=JSON&IOSRemind=true&PushType=NOTICE&RegionId=cn-hangzhou&Signature=MEbFt3SwEHqtLIMHsvwdx9fO5hM%3D&SignatureMethod=HMAC-SHA1&SignatureNonce=cpcU3TWloqrJsnxK2K3Fb19ohR3dHcVB&SignatureVersion=1.0&Summary=hello&Target=DEVICE&TargetValue=b035d8979d5b4c8db8482d49b1a2f0e7%2C6fbf8c1c023f4fc5819892ce66d6777e%2C9e4f23fa76664cd994b7da8859ecfa85%2C539f8d9e83a148b097807f9e6783a133%2C8341bcc3bf7f45728d9b21f3efc509ce%2Ca25f041394454af9bf31b14e87a54ba7%2C6159ef75c04047968da8b13c65b6fd83%2C80b05deddc7a4b7c95364039c9dbfcdf%2C472a31b69b8145ffb76bf17aa52b322c&Timestamp=2018-02-12T07%3A49%3A40Z&Version=2016-08-01&iOSApnsEnv=PRODUCT&iOSBadge=1&iOSExtParameters=%7B%22c%22%3A%22%22%2C%22id%22%3A%22220024%22%2C%22new_morning_paper%22%3A%221%22%2C%22schema%22%3A%22cailianshe%3A%2F%2Ftelegram_detail%3Fdetail_id%3D220024%22%2C%22title%22%3A%22%E8%B4%A2%E8%81%94%E7%A4%BE%22%2C%22topicid%22%3A%22%22%2C%22type%22%3A%22-1%22%7D&iOSRemind=true"
	values, err := url.ParseQuery(str)
	if err != nil {
		panic(err)
	}
	jsonData, err := json.Marshal(values)
	if err != nil {
		panic(err)
	}
	log = log.With("pushData", string(jsonData))
	log.Debugf("")
	// i := 0
	// for i < 1000 {
	// 	i++
	// 	aliLog.Debugf("log %d", i)
	// }

	time.Sleep(5 * time.Second)
}
