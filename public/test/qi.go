package main

import (
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
	"log"
)

func main() {
	log.Println(GetQiniuUploadToken())
}

var ak = "HTTyjWkdHISJbKTD0n3OZ_2UPt-AvBKdPRZs2wxQ"
var sk = "9Xp6-AlBO9mqP9iyPKsYgxVadj93sIEcfdGxxnG9"

func GetQiniuUploadToken() (token string) {
	putPolicy := storage.PutPolicy{
		Scope: "hobby-box",
	}
	putPolicy.Expires = 7200 //示例2小时有效期
	mac := qbox.NewMac(ak, sk)
	return putPolicy.UploadToken(mac)

}
