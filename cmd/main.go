package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/ronething/mp-dev/config"
	"github.com/ronething/mp-dev/storage"

	"github.com/ronething/mp-dev/server"
)

var (
	filePath string // 配置文件路径
	help     bool   // 帮助
)

func usage() {
	fmt.Fprintf(os.Stdout, `wechat-mp - simlpe wechat mp handler
Usage: wechat-mp [-h help] [-c ./config.yaml]                                    
Options:                                    
`)
	flag.PrintDefaults()
}

func main() {
	flag.StringVar(&filePath, "c", "./config.yaml", "配置文件所在")
	flag.BoolVar(&help, "h", false, "帮助")
	flag.Usage = usage
	flag.Parse()
	if help {
		flag.PrintDefaults()
		return
	}

	// 设置配置文件和静态变量
	config.SetConfig(filePath)

	// 初始化
	storage.InitWechatConfig()
	storage.InitThirdService()

	e := server.CreateEngine()
	host := config.Config.GetString("server.host")
	if err := e.Start(host); err != nil {
		fmt.Printf("启动服务失败, err:%v\n", err)
		return
	}
}
