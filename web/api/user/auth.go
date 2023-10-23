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

type wechatLoginVerify struct {
	JsToken string `binding:"required" json:"js_token"`
}

type WeChatRegisterVerify struct {
	Account      string `binding:"required"`
	Avatar       string
	TelPhone     string
	Mail         string
	Username     string //昵称
	WechatOpenID string `binding:"required" json:"wechat_open_id"`
}

var JwtSignKey = []byte("fantuanKey")

func makeLogin(user *model.User) (JT string, err error) {
	userClaims := CommonClaims{
		Id:       user.Id,
		Username: user.Username,
		Account:  user.Account,
		Telphone: user.TelPhone,
		StandardClaims: jwt.StandardClaims{
			NotBefore: time.Now().Unix() - 60, //生效时间
			ExpiresAt: time.Now().Unix() + 86400*3,
			Issuer:    "UserCenter", //签发人
		},
	}

	//登录成功，记录登录信息
	user.LastLoginTime = time.Now().Unix()
	userJwt := jwt.NewWithClaims(jwt.SigningMethodHS256, userClaims)

	jwtToken, err := userJwt.SignedString(JwtSignKey)
	if err != nil {
		return "", errors.New("加密SignKey错误，登录失败")
	}
	user.LastToken = jwtToken
	user.LastTokenExpire = time.Now().Unix() + 3600*2
	bootstrap.DB.Save(user)

	return jwtToken, nil
}

func makeRegister(u *model.User) (*model.User, error) {
	user := model.User{
		Account:      u.Account,
		Avatar:       u.Avatar,
		TelPhone:     u.TelPhone,
		Mail:         u.Mail,
		Username:     u.Username,
		Password:     u.Password,
		WechatOpenID: u.WechatOpenID,
		CreatedAt:    time.Now().Unix(),
		ShortDomain:  u.ShortDomain,
	}
	result := bootstrap.DB.Create(&user)
	if result.Error != nil {
		log.Printf("有用户注册失败，原因:%s。参数%+v", result.Error, user)
		return nil, errors.New(fmt.Sprintf(api.FailedMsg, "系统内部错误，暂时无法注册"))
	}
	return &user, nil
}
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
	if res.Error != nil {
		api.Base.Failed(ctx, "服务器内部错误暂时无法登录auth-50001")
		return
	}
	//加盐
	mdPass := tools.HexMd5(loginStruct.Password, "HexFanTuan")
	if mdPass != user.Password {
		api.Base.Failed(ctx, fmt.Sprintf(api.FailedMsg, "密码输入错误"))
		return
	}

	jwtToken, err := makeLogin(user)
	if err != nil {
		api.Base.Failed(ctx, err.Error())
		return
	}
	ret := make(map[string]string)
	ret["jwt_token"] = jwtToken

	api.Base.Success(ctx, ret)
}

func LoginWithWechat(ctx *gin.Context) {
	var binding wechatLoginVerify
	if err := ctx.ShouldBindJSON(&binding); err != nil {
		api.Base.Failed(ctx, "缺少必要参数")
		return
	}

	wechatResp, err := bootstrap.Wechat.Login(binding.JsToken)
	if err != nil {
		api.Base.Failed(ctx, "服务器暂时繁忙，授权登录失败。")
		return
	}
	if wechatResp.ErrCode != 0 {
		errData := make(map[string]int)
		errData["wechat_err_code"] = wechatResp.ErrCode
		api.Base.FailedWithContent(ctx, "微信授权登录失败", errData)
		return
	}
	openId := wechatResp.OpenID
	sessionKey := wechatResp.SessionKey
	//var openId = "xxx"
	//var sessionKey = "xxx"
	//var err error

	//查询openID是否注册过，如果没有返回前端询问是否注册新号
	user := new(model.User)
	res := bootstrap.DB.Where("wechat_open_id = ?", openId).First(user)

	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		//如果未注册用户，则生成一个新号
		user.WechatOpenID = openId
		user.WechatSessionKey = sessionKey
		user.Password = "wechat-register"
		user.Username = "微信用户" + openId
		user.Account = "微信用户" + openId
		user.Avatar = "/default_avatar.jpg"
		user.ShortDomain = ""
		user, err = makeRegister(user)
		if err != nil {
			api.Base.Failed(ctx, err.Error())
			return
		}
	}
	if res.Error != nil && !errors.Is(res.Error, gorm.ErrRecordNotFound) {
		api.Base.Failed(ctx, "服务器内部错误暂时无法登录auth-50002 "+res.Error.Error())
		return
	}

	jwtToken, err := makeLogin(user)
	if err != nil {
		api.Base.Failed(ctx, err.Error())
		return
	}
	ret := make(map[string]string)
	ret["jwt_token"] = jwtToken
	ret["open_id"] = openId
	ret["session_key"] = sessionKey

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
			return nilJ, errors.New("登录状态过期")
		}
		log.Printf("解析JWT错误:%s,jwt-string:%s\n", ve.Error(), authorization)

		return nilJ, err
	}
	jwtInfo := *loginStatus.Claims.(*CommonClaims)
	return jwtInfo, nil
}

func Mine(ctx *gin.Context) {
	// 获取中间件传过来的user
	u, _ := ctx.Get("user")
	user := u.(*model.User)
	//解析auth中的authorization
	a := ctx.GetHeader("Authorization")
	authorization, _ := ParseJWT(a)
	if authorization.ExpiresAt-time.Now().Unix() < 3600*6 {
		// token即将过期，生成一个新的jwt token
		newJwt, err := makeLogin(user)
		if err != nil {
			api.Base.Failed(ctx, "获取个人信息或更新登录状态产生未知错误。")
			return
		}
		user.LastToken = newJwt
	}

	api.Base.Success(ctx, user)
	return
}
