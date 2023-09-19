package config

import (
	"bufio"
	"io"
	"log"
	"os"
)

func Load() {
	iniFile, err := os.Open("./config/dev.ini")
	if err != nil {
		log.Panicln("不能打开INI文件", err)
	}
	//tmpRead := make([]byte, 9)
	defer iniFile.Close()
	reader := bufio.NewReader(iniFile)
	for {
		line, _, err := reader.ReadLine()

		if err == io.EOF {
			//log.Println(string(tmpRead[:n]))
			log.Println("配置加载结束")
			break
		}
		if err != nil {
			log.Panicln("读取错误")
			break
		}
		log.Println(string(line))
	}

}
