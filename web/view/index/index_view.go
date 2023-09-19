package web_view

import "github.com/gin-gonic/gin"

func Page(ctx *gin.Context) {
	ctx.HTML(200, "test", "aa")
}
