package initialize

import (
	"ganxue-server/static"
	"ganxue-server/utils/db"

	"github.com/gangantongxue/ggl"
)

func InitAll() {
	// 加载配置文件
	if err := LoadConfig(); err != nil {
		ggl.Error("加载配置文件失败", ggl.Err(err.ToError()))
	}

	logCfg := ggl.DefaultConfig()
	logCfg.LogFileName = "ganxue_log_2006-01-02.log"
	logCfg.LogFileDir = "./log"
	logCfg.ToConsole = true
	ggl.New(logCfg)

	// 初始化数据库
	db.Init()

	// 初始化静态文件
	static.Init()
}

func CloseAll() {
	db.Close()
}
