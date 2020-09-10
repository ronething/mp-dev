package config

import (
	"github.com/fsnotify/fsnotify"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var Config *viper.Viper

func SetConfig(filePath string) {
	log.Infof("[config] run the env with:%s", filePath)

	Config = viper.New()
	Config.SetConfigFile(filePath)
	if err := Config.ReadInConfig(); err != nil {
		log.Fatalf("[config] read config err: %v", err)
	}

	// set log by default
	setLog()

	watchFileConfig()
}

// set log level
func setLog() {
	l := Config.Get("log.level")

	if l == "" {
		log.SetLevel(log.InfoLevel)
	} else if l == "debug" {
		log.SetLevel(log.DebugLevel)
	} else if l == "error" {
		log.SetLevel(log.ErrorLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}

}

//watchFileConfig 监听文件变化
func watchFileConfig() {
	Config.WatchConfig()
	Config.OnConfigChange(func(e fsnotify.Event) {
		log.Warnf("config file change: %v %v", e.Name, e.Op)
		setLog()
	})
}
