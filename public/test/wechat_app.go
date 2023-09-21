package main

import (
	"github.com/medivhzhan/weapp/v3"
	"log"
)

func main() {
	sdk := weapp.NewClient("wx5011e6982b77cce5", "dbe4b4806b042e5ec7fba5fe809acffd")
	resp, err := sdk.Login("0e3mKWZv3sJwo13Zoo0w3vzoLr1mKWZx")
	if err != nil {
		log.Println("err:", err)
		return
	}
	log.Printf("%+v\n", resp)
}
