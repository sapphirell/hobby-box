package api

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type BaseApi struct{}

var Base BaseApi

var (
	SuccessMsg = "提交成功啦！\U0001F97A"
	FailedMsg  = "操作失败啦！\U0001F97A %s"
)

func (b *BaseApi) Success(ctx *gin.Context, data interface{}) {
	ctx.Header("Content-Type", "application/json")
	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"msg":    SuccessMsg,
		"data":   data,
	})
}
func (b *BaseApi) Failed(ctx *gin.Context, msg string) {
	ctx.Header("Content-Type", "application/json")
	ctx.JSON(http.StatusOK, gin.H{
		"status": "failed",
		"msg":    msg,
		"data":   nil,
	})
}

func (b *BaseApi) FailedWithContent(ctx *gin.Context, msg string, content interface{}) {
	ctx.Header("Content-Type", "application/json")
	ctx.JSON(http.StatusOK, gin.H{
		"status": "failed",
		"msg":    msg,
		"data":   content,
	})
}

// GetParam 获取JSON参数
func (b *BaseApi) GetParam(ctx *gin.Context, key string) interface{} {
	j := make(map[string]interface{})
	err := ctx.ShouldBindJSON(&j)
	if err != nil {
		log.Println("无法解析JSON参数", err)
	}
	if key == "" {
		return j
	}
	return j[key]
}
