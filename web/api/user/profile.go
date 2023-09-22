package user

import (
	"github.com/gin-gonic/gin"
	"sukitime.com/v2/web/api"
)

type UpdateProfile struct {
	Id       string `binding:"required"`
	Avatar   string
	TelPhone string
	Mail     string
	Username string //昵称
}

func UpdateProfileFn(ctx *gin.Context) {
	var binding UpdateProfile
	if err := ctx.ShouldBindJSON(&binding); err != nil {
		api.Base.Failed(ctx, "缺少必要参数")
		return
	}

}
