package bootstrap

import "github.com/medivhzhan/weapp/v3"

func InitSocial() {
	//初始化微信APP盒子
	Wechat = weapp.NewClient("wx5011e6982b77cce5", "dbe4b4806b042e5ec7fba5fe809acffd")
}
