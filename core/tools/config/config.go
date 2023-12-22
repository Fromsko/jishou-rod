package config

import (
	"fmt"
	"os"

	"github.com/Fromsko/gouitls/logs"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	Log       *logrus.Logger
	Conifg    *viper.Viper
	FieldList = []string{
		"WeatherKey", "TemplateID",
		"CnameData", "CnameImage",
		"AppID", "AppSecret",
		"UserName", "PassWord",
	}
)

func init() {
	Log = logs.InitLogger()
	Conifg = InitConfig()
}

func InitConfig() *viper.Viper {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	if err := initConfig(); err != nil {
		return nil
	}

	return viper.GetViper()
}

func initConfig() error {
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return initConfigFile()
		}
		return err
	}
	defer checkConfig()
	return nil
}

func initConfigFile() error {
	for _, v := range FieldList {
		viper.SetDefault(v, "")
	}

	if err := viper.WriteConfigAs("config.yaml"); err != nil {
		return err
	} else {
		fmt.Println("请填写相应的字段值")
		os.Exit(0)
	}
	return nil
}

func checkConfig() {
	var missingFields []string

	for _, field := range FieldList {
		if !viper.IsSet(field) || viper.GetString(field) == "" {
			missingFields = append(missingFields, field)
		}
	}

	if len(missingFields) > 0 {
		fmt.Println("配置文件缺少必要字段:", missingFields)
		os.Exit(0)
	}
}
