package tools

import (
	"context"
	"crypto/md5"
	"errors"
	"fmt"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"path/filepath"
	"strconv"
	"time"
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

// Download 通用下载方法
func Download(url string, timeout time.Duration) (store string, fileName string, err error) {
	downloadChan := make(chan string, 1)
	timeoutChan := make(chan string, 1)
	errChan := make(chan error, 1)
	//建议存储的fileName
	fileName = HexMd5(url, strconv.Itoa(rand.Int())) + ".jpg"
	go func() {
		time.Sleep(timeout)
		timeoutChan <- "*"
	}()

	go func() {
		i, err := http.Get(url)
		defer i.Body.Close()
		if err != nil {
			errChan <- err
			return
		}

		content, err := ioutil.ReadAll(i.Body)
		if err != nil {
			errChan <- err
			return
		}
		storePath, _ := filepath.Abs(fmt.Sprintf("./public/tmp/red_book/%s", fileName))
		err = ioutil.WriteFile(storePath, content, 0666)
		if err != nil {
			errChan <- err
			return
		}
		downloadChan <- storePath
	}()
	select {
	case v := <-downloadChan:
		return v, fileName, nil
	case err := <-errChan:
		return "", fileName, err
	case <-timeoutChan:
		return "", fileName, errors.New(fmt.Sprintf("下载图片%s超时", url))

	}
}

func Upload2QiNiu(url string, save string) (savePath string, err error) {
	putPolicy := storage.PutPolicy{
		Scope: "hobby-box",
	}
	mac := qbox.NewMac("HTTyjWkdHISJbKTD0n3OZ_2UPt-AvBKdPRZs2wxQ",
		"9Xp6-AlBO9mqP9iyPKsYgxVadj93sIEcfdGxxnG9")
	upToken := putPolicy.UploadToken(mac)
	cfg := storage.Config{}
	cfg.Region = &storage.ZoneHuadongZheJiang2
	cfg.UseHTTPS = true
	cfg.UseCdnDomains = false
	// 构建表单上传的对象
	formUploader := storage.NewFormUploader(&cfg)
	ret := storage.PutRet{}
	err = formUploader.PutFile(context.Background(), &ret, upToken, save, url, nil)
	if err != nil {
		return "", err
	}
	//fmt.Println(ret.Key, ret.Hash)

	return ret.Key, nil
}
