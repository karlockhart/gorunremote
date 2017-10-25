package config

import (
	"github.com/Sirupsen/logrus"
	"github.com/spf13/viper"
)

func LoadConfig() {
	viper.SetConfigName("gorunremote")
	viper.AddConfigPath(".")
	viper.AddConfigPath("~/.gorunremote")
	viper.AddConfigPath("./etc")
	viper.AddConfigPath("/etc/gorunremote/")
	logrus.Info("Loading Config")
	e := viper.ReadInConfig()
	if e != nil {
		logrus.Fatal(e)
	}
}
