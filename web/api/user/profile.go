package user

import (
	"github.com/gin-gonic/gin"
	"sukitime.com/v2/bootstrap"
	"sukitime.com/v2/web/api"
	"sukitime.com/v2/web/model"
)

type UpdateProfile struct {
	Avatar   string
	TelPhone string
	Username string //昵称
}

type UpdateLoginAccountVerify struct {
	Account string `json:"account" binding:"required"`
}

func UpdateProfileFn(ctx *gin.Context) {
	var binding UpdateProfile
	if err := ctx.ShouldBindJSON(&binding); err != nil {
		api.Base.Failed(ctx, "缺少必要参数")
		return
	}
	u, _ := ctx.Get("user")
	user := u.(model.User)
	if binding.TelPhone != user.TelPhone {
		user.TelPhone = binding.TelPhone
	}
	if binding.Avatar != user.Avatar {
		user.Avatar = binding.Avatar
	}
	if binding.Username != user.Username {
		user.Username = binding.Username
	}
	bootstrap.DB.Save(user)

	api.Base.Success(ctx, "")
}

func UpdateLoginAccount(ctx *gin.Context) {
	var binding UpdateLoginAccountVerify
	if err := ctx.ShouldBindJSON(&binding); err != nil {
		api.Base.Failed(ctx, "缺少必要参数")
		return
	}
}
