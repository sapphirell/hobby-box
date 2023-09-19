package web

import (
	gin "github.com/gin-gonic/gin"
	"log"
	"net/http"
	"path/filepath"
	"sukitime.com/v2/web/api"
	"sukitime.com/v2/web/api/box"
	"sukitime.com/v2/web/api/spider"
	User "sukitime.com/v2/web/api/user"
	"sukitime.com/v2/web/model"
	IndexView "sukitime.com/v2/web/view/index"
)

func LoadRouter() {
	router := gin.Default()
	router.GET("/", func(context *gin.Context) {
		context.String(http.StatusOK, "api")
	})
	router.GET("/index", IndexView.Page)
	router.POST("/login", User.Login)
	router.POST("/red_book_spider", spider.GetDownloadList)

	needLoginGroup := router.Group("/with-state")
	{
		needLoginGroup.Use(LoginVerify())
		needLoginGroup.GET("/mine", User.Mine)
		needLoginGroup.POST("/register", User.Register)
		needLoginGroup.GET("/box-items", box.ItemList)
		needLoginGroup.POST("/add-items", box.AddItems)
	}

	go func() {
		//启动tls
		certPath, _ := filepath.Abs("./web/cert/api.fantuanpu.com/api.fantuanpu.com_bundle.pem")
		keyPath, _ := filepath.Abs("./web/cert/api.fantuanpu.com/api.fantuanpu.com.key")
		err := router.RunTLS(":443",
			certPath,
			keyPath)
		if err != nil {
			log.Fatal("https api启动错误", err)
			return
		}

	}()

	go func() {
		//启动http
		err := router.Run(":8080")
		if err != nil {
			log.Fatal("http启动错误")
			return
		}
	}()

}

func LoginVerify() gin.HandlerFunc {
	return func(c *gin.Context) {
		authorization := c.GetHeader("Authorization")
		if authorization == "" {
			api.Base.Failed(c, "missing header params")
			c.Abort()
			return
		}
		usr, err := User.ParseJWT(authorization)
		if err != nil {
			api.Base.Failed(c, "parse jwt-token failed")
			c.Abort()
			return
		}
		userInfo, err := model.UserModel.GetUserInfoById(usr.Id)
		if err != nil {
			api.Base.Failed(c, "User not fund")
			c.Abort()
			return
		}
		c.Set("user", userInfo)
		log.Printf("userInfo:%+v", userInfo)
		c.Next()
	}
}
