package box

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"strconv"
	"sukitime.com/v2/web/api"
	"sukitime.com/v2/web/model"
)

type addItemsVerify struct {
	Name        string  `binding:"required"`
	Image       string  `binding:"required"`
	BuyTime     int64   `binding:"required" json:"buy_time"`
	Price       float64 `binding:"required"`
	Status      int64   `binding:"required"`
	NextPayTime int64   `binding:"required" json:"next_pay_time"`
	Type        string  `binding:"required"`
}

type updateItemsVerify struct {
	addItemsVerify
	Id string `binding:"required" json:"id"`
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
	i.Uid = user.Id
	i.Name = binding.Name
	i.Image = binding.Image
	i.BuyTime = binding.BuyTime
	i.Price = binding.Price
	i.Status = binding.Status
	i.NextPayTime = binding.NextPayTime
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
	id, err := strconv.Atoi(binding.Id)
	if err != nil {
		api.Base.Failed(ctx, "params invalid: ID must be int")
		return
	}
	u, _ := ctx.Get("user")
	user := u.(*model.User)

	items, err := model.ItemsModel.GetItemsInfoById(int64(id))
	if err != nil {
		api.Base.Failed(ctx, fmt.Sprintf(api.FailedMsg, "无法查询到对应数据"))
		return
	}
	if items.Uid != user.Id {
		api.Base.Failed(ctx, fmt.Sprintf(api.FailedMsg, "不可访问的物品"))
		return
	}

	api.Base.Success(ctx, "")
}
