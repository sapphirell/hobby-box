package model

import (
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
	"log"
	"sukitime.com/v2/bootstrap"
)

// Items 我的物品表
type Items struct {
	Id          int64   `json:"id"`
	Name        string  `json:"name"`
	Image       string  `json:"image"`
	Uid         int64   `json:"uid"`
	BuyTime     int64   `json:"buy_time"`
	Price       float64 `json:"price"`
	Status      int64   `json:"status"`
	NextPayTime int64   `json:"next_pay_time"`
	Type        string  `json:"type"`
}

var ItemsModel Items

func (Items) TableName() string {
	return "items"
}

func (Items) GetItemsInfoById(id int64) (*Items, error) {
	items := new(Items)
	res := bootstrap.DB.Where("id = ?", id).First(items)
	return items, res.Error
}

func (Items) GetMyItemsList(uid int64, theType string, page int) ([]*Items, error) {
	var items []*Items
	var res *gorm.DB
	pageSize := 20
	page--

	cacheKey := fmt.Sprintf(ItemListCacheKey, theType, page, uid)
	redisRes, err := bootstrap.RedisClient.Get(cacheKey).Result()
	if err != nil {
		log.Printf("无法从Redis获取数据%s", err.Error())
	}
	err = json.Unmarshal([]byte(redisRes), &items)
	if err != nil {
		log.Printf("解析RedisJson错误:%s", err.Error())
	}

	if theType == "" {
		res = bootstrap.DB.Where("uid = ?", uid).Limit(pageSize).Offset(pageSize * page).Find(&items)
	} else {
		res = bootstrap.DB.Where("uid = ? and type = ?", uid, theType).Limit(pageSize).Offset(pageSize * page).Find(&items)
	}

	return items, res.Error
}

// AddItems 添加我的物品并删除缓存
func (Items) AddItems(i *Items, flushCache bool) {
	bootstrap.DB.Create(i)
	//分类信息列表缓存
	typeCacheKey := fmt.Sprintf(ItemListCacheKey, i.Type, 1, i.Uid)
	//总列表缓存
	listCacheKey := fmt.Sprintf(ItemListCacheKey, "", 1, i.Uid)
	if flushCache {
		bootstrap.RedisClient.Del(listCacheKey)
		bootstrap.RedisClient.Del(typeCacheKey)
	}
}
