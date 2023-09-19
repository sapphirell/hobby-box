package main

import (
	"golang.org/x/sync/singleflight"
	"log"
)

func GetDataFromDB() string {
	return "GetForDB"
}

// 使用single_flight实现防止缓存击穿
func main() {
	var sf singleflight.Group

	res, err, _ := sf.Do("MissionKey1", func() (interface{}, error) {
		return GetDataFromDB(), nil
	})
	if err != nil {
		log.Println(err)
	}
	log.Println(res)
}
