package tools

import (
	"crypto/md5"
	"fmt"
	"log"
)

// HexMd5 计算32位md5， salt可选加盐
func HexMd5(s string, salt string) string {
	b := []byte(s)
	if salt != "" {
		b = append(b, []byte(salt)...)
	}
	hash := md5.New()
	_, err := hash.Write(b)
	if err != nil {
		log.Println("生成md5错误", err)
	}
	sum := hash.Sum(nil)
	return fmt.Sprintf("%x", sum)
}
