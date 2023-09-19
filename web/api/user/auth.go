package user

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"sukitime.com/v2/bootstrap"
	"sukitime.com/v2/tools"
	"sukitime.com/v2/web/api"
	"sukitime.com/v2/web/model"
	"time"
)

type Auth struct {
	BaseApi *api.BaseApi
}
type Message struct {
}

// CommonClaims jwtClaims
type CommonClaims struct {
	Id       int64  `json:"id"`
	Username string `json:"username"`
	Account  string `json:"account"`
	Telphone string `json:"telphone"`
	jwt.StandardClaims
}

type loginVerify struct {
	Account  string `binding:"required"`
	Password string `binding:"required"`
}

type RegisterVerify struct {
	Account  string `binding:"required"`
	Password string `binding:"required"`
	Avatar   string
	TelPhone string
	Mail     string
	Username string //昵称
}

var JwtSignKey = []byte("fantuanKey")

func Login(ctx *gin.Context) {
	loginStruct := loginVerify{}
	if err := ctx.ShouldBind(&loginStruct); err != nil {
		log.Println("错误", err.Error())
		api.Base.Failed(ctx, fmt.Sprintf(api.FailedMsg, "登录信息好像输错了"))
		return
	}
	var user *model.User
	user = new(model.User)

	res := bootstrap.DB.Where("account = ?", loginStruct.Account).First(user)
	log.Println(res.Error)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		api.Base.Failed(ctx, fmt.Sprintf(api.FailedMsg, "您尚未注册账号"))
		return
	}
	//加盐
	mdPass := tools.HexMd5(loginStruct.Password, "HexFanTuan")
	if mdPass != user.Password {
		api.Base.Failed(ctx, fmt.Sprintf(api.FailedMsg, "密码输入错误"))
		return
	}

	userClaims := CommonClaims{
		Id:       user.Id,
		Username: user.Username,
		Account:  user.Account,
		Telphone: user.TelPhone,
		StandardClaims: jwt.StandardClaims{
			NotBefore: time.Now().Unix() - 60, //生效时间
			ExpiresAt: time.Now().Unix() + 3600,
			Issuer:    "UserCenter", //签发人
		},
	}

	//登录成功，记录登录信息
	user.LastLoginTime = time.Now().Unix()
	userJwt := jwt.NewWithClaims(jwt.SigningMethodHS256, userClaims)

	jwtToken, err := userJwt.SignedString(JwtSignKey)
	if err != nil {
		log.Println("加密signKey错误", err)
		api.Base.Failed(ctx, fmt.Sprintf(api.FailedMsg, "登录服务器错误"))
		return
	}
	user.LastToken = jwtToken
	user.LastTokenExpire = time.Now().Unix() + 3600*2
	bootstrap.DB.Save(user)

	ret := make(map[string]string)
	ret["jwt_token"] = jwtToken

	api.Base.Success(ctx, ret)
}

func ParseJWT(authorization string) (CommonClaims, error) {
	loginStatus, err := jwt.ParseWithClaims(authorization, &CommonClaims{}, func(token *jwt.Token) (interface{}, error) {
		return JwtSignKey, nil
	})
	nilJ := *new(CommonClaims)
	if err != nil {
		ve, _ := err.(*jwt.ValidationError)
		if ve.Errors == jwt.ValidationErrorExpired {
			// 登录超时
			return nilJ, err
		}
		log.Printf("解析JWT错误:%s,jwt-string:%s\n", ve.Error(), authorization)

		return nilJ, err
	}
	jwtInfo := *loginStatus.Claims.(*CommonClaims)
	return jwtInfo, nil
}

func Mine(ctx *gin.Context) {
	Authorization := ctx.GetHeader("Authorization")
	userStatus, err := ParseJWT(Authorization)
	if err != nil {
		api.Base.Failed(ctx, err.Error())
		return
	}
	api.Base.Success(ctx, userStatus)
	return
}

func Register(ctx *gin.Context) {
	registerVerify := RegisterVerify{}
	if err := ctx.ShouldBind(&registerVerify); err != nil {
		log.Println("注册参数错误", err.Error())
		api.Base.Failed(ctx, fmt.Sprintf(api.FailedMsg, "注册信息填写错了！"))
		return
	}

	//填充默认参数
	if registerVerify.Avatar == "" {
		registerVerify.Avatar = "default_avatar.jpg"
	}
	if registerVerify.Username == "" {
		registerVerify.Username = registerVerify.Account
	}
	user := model.User{
		Account:   registerVerify.Account,
		Avatar:    registerVerify.Avatar,
		TelPhone:  registerVerify.TelPhone,
		Mail:      registerVerify.Mail,
		Username:  registerVerify.Username,
		Password:  registerVerify.Password,
		CreatedAt: time.Now().Unix(),
	}
	result := bootstrap.DB.Create(&user)
	if result.Error != nil {
		api.Base.Failed(ctx, fmt.Sprintf(api.FailedMsg, "系统内部错误，暂时无法注册。"))
		log.Printf("有用户注册失败，原因:%s。参数%+v", result.Error, user)
		return
	}
	api.Base.Success(ctx, "")
}
