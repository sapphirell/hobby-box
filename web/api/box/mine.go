package box

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"strconv"
	"sukitime.com/v2/web/api"
	"sukitime.com/v2/web/model"
	"time"
)

type addItemsVerify struct {
	Name        string  `binding:"required"`
	Image       string  `binding:"required"`
	BuyTime     string  `binding:"required" json:"buy_time"`
	Price       float64 `binding:"required"`
	Status      int64   `binding:"required"`
	NextPayTime string  `binding:"required" json:"next_pay_time"`
	Type        string  `binding:"required"`
}

type updateItemsVerify struct {
	addItemsVerify
	Id int64 `binding:"required" json:"id"`
}

func ItemList(ctx *gin.Context) {
	u, _ := ctx.Get("user")
	p := ctx.DefaultQuery("page", "1")
	theType := ctx.DefaultQuery("type", "")
	page, err := strconv.Atoi(p)
	if err != nil {
		api.Base.Failed(ctx, "params `page` invalid")
	}
	user := u.(*model.User)
	list, err := model.ItemsModel.GetMyItemsList(user.Id, theType, page)
	if err != nil {
		log.Printf("查询uid:%d,返回异常%s\n", user.Id, err.Error())
		api.Base.Failed(ctx, "not find list")
		return
	}
	api.Base.Success(ctx, list)
}

func AddItems(ctx *gin.Context) {
	var i model.Items
	var binding addItemsVerify
	u, _ := ctx.Get("user")
	user := u.(*model.User)

	if err := ctx.ShouldBindJSON(&binding); err != nil {
		api.Base.Failed(ctx, "verify failed:"+err.Error())
		return
	}
	buyTime, err := time.Parse("2006-01-02", binding.BuyTime)
	if err != nil {
		log.Println("时间格式化失败:将", binding.BuyTime, "格式化为时间戳")
	}
	nextPayTime, err := time.Parse("2006-01-02 15:04", binding.NextPayTime)
	if err != nil {
		log.Println("时间格式化失败:将", binding.BuyTime, "格式化为时间戳")
	}
	i.Uid = user.Id
	i.Name = binding.Name
	i.Image = binding.Image
	i.BuyTime = buyTime.Unix()
	i.Price = binding.Price
	i.Status = binding.Status
	i.NextPayTime = nextPayTime.Unix()
	i.Type = binding.Type

	model.ItemsModel.AddItems(&i, true)
	api.Base.Success(ctx, "")
}

func UpdateItems(ctx *gin.Context) {
	var binding updateItemsVerify
	if err := ctx.ShouldBindJSON(&binding); err != nil {
		api.Base.Failed(ctx, "params invalid: "+err.Error())
		return
	}
	u, _ := ctx.Get("user")
	user := u.(*model.User)

	items, err := model.ItemsModel.GetItemsInfoById(binding.Id)
	if err != nil {
		api.Base.Failed(ctx, fmt.Sprintf(api.FailedMsg, "无法查询到对应数据"))
		return
	}
	if items.Uid != user.Id {
		api.Base.Failed(ctx, fmt.Sprintf(api.FailedMsg, "不可访问的物品"))
		return
	}

	buyTime, err := time.Parse("2006-01-02", binding.BuyTime)
	if err != nil {
		log.Println("时间格式化失败:将", binding.BuyTime, "格式化为时间戳")
	}
	nextPayTime, err := time.Parse("2006-01-02 15:04", binding.NextPayTime)
	if err != nil {
		log.Println("时间格式化失败:将", binding.NextPayTime, "格式化为时间戳")
	}
	items.Name = binding.Name
	items.Image = binding.Image
	items.BuyTime = buyTime.Unix()
	items.Price = binding.Price
	items.Status = binding.Status
	items.NextPayTime = nextPayTime.Unix()
	items.Type = binding.Type

	model.ItemsModel.UpdateItem(items)

	api.Base.Success(ctx, "")
}
