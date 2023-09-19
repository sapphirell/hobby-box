package model

import "sukitime.com/v2/bootstrap"

// User 用户表
type User struct {
	Id              int64  `json:"id"`
	LastLoginTime   int64  `json:"last_login_time"`
	Account         string `json:"account"`
	Avatar          string `json:"avatar"`
	TelPhone        string `json:"tel_phone"`
	Mail            string `json:"mail"`
	LastToken       string `json:"last_token"`
	LastTokenExpire int64  `json:"last_token_expire"`
	Username        string `json:"username"`
	Password        string `json:"password"`
	CreatedAt       int64  `json:"created_at"`
}

var UserModel User

func (User) TableName() string {
	return "user"
}

func (User) GetUserInfoById(id int64) (*User, error) {
	usr := new(User)
	res := bootstrap.DB.Where("id = ?", id).First(usr)
	return usr, res.Error
}
