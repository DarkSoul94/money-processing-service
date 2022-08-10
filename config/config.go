package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Release  bool
	AppName  string
	HTTPport string
	LogDir   string
	LogFile  string
	DbHost   string
	DbPort   string
	DBLogin  string
	DBPass   string
	DbName   string
	DBArgs   string
}

// InitConfig - load config from config.yml
func InitConfig() Config {
	viper.AddConfigPath("./config")
	viper.SetConfigName("config")

	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	conf := Config{
		Release:  viper.GetBool("release"),
		AppName:  viper.GetString("name"),
		HTTPport: viper.GetString("http_port"),
		LogDir:   viper.GetString("log.dir"),
		LogFile:  viper.GetString("log.file"),
		DbHost:   viper.GetString("db.host"),
		DbPort:   viper.GetString("db.port"),
		DBLogin:  viper.GetString("db.login"),
		DBPass:   viper.GetString("db.pass"),
		DbName:   viper.GetString("db.name"),
		DBArgs:   viper.GetString("db.args"),
	}

	return conf
}
