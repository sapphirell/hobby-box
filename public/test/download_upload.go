package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"sukitime.com/v2/tools"
	"time"
)

func main() {
	dArr := []string{
		"https://sns-img-bd.xhscdn.com/6ed83d93-353e-d023-3109-485b26af3409?imageView2/2/w/1920/format/jpg|imageMogr2/strip",
		"https://sns-img-bd.xhscdn.com/7dac19bf-36c1-0cb6-7d90-acd23eac4bc4?imageView2/2/w/1920/format/jpg|imageMogr2/strip",
	}

	store, save, err := download(dArr[0], 3*time.Second)
	if err != nil {
		log.Println("下载失败", err)
		return
	}
	log.Println("存储目录", store)
	save = fmt.Sprintf("sp/%s/%s", time.Now().Format("2006_01_02"), save)
	upload2QiNiu(store, save)
	//for _, s := range dArr {
	//	log.Println("download:", s)
	//	store, err := download(s, 1*time.Second)
	//	if err != nil {
	//		log.Println("下载失败", err)
	//	}
	//	log.Println("存储目录", store)
	//}
}

func download(url string, timeout time.Duration) (store string, fileName string, err error) {
	downloadChan := make(chan string, 1)
	timeoutChan := make(chan string, 1)
	errChan := make(chan error, 1)
	//建议存储的fileName
	fileName = tools.HexMd5(url, "") + ".jpg"
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
		storePath, _ := filepath.Abs(fmt.Sprintf("../tmp/red_book/%s", fileName))
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

func upload2QiNiu(url string, save string) {
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
	err := formUploader.PutFile(context.Background(), &ret, upToken, save, url, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(ret.Key, ret.Hash)
}
