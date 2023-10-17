package box

import (
	"github.com/gin-gonic/gin"
	"math/rand"
	"strconv"
	"sukitime.com/v2/tools"
	"sukitime.com/v2/web/api"
	"sukitime.com/v2/web/model"
	"time"
)

func QiniuToken(ctx *gin.Context) {
	get, exists := ctx.Get("user")
	if !exists {
		api.Base.Failed(ctx, "登录状态异常")
		return
	}
	user := get.(*model.User)
	path := "/box/" + strconv.FormatInt(user.Id, 10) + "/"
	path = path + tools.HexMd5(time.Now().String(), strconv.FormatFloat(rand.Float64(), 'f', 2, 64))
	token := tools.GetQiniuUploadToken(path)
	ret := make(map[string]string)
	ret["token"] = token
	ret["path"] = path

	api.Base.Success(ctx, ret)
}
