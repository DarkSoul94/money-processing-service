package config

import (
	"os"

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
		DBLogin:  viper.GetString("db.login"),
		DBPass:   viper.GetString("db.pass"),
		DbName:   viper.GetString("db.name"),
	}

	if dbHost := os.Getenv("DBHost"); len(dbHost) != 0 {
		conf.DbHost = dbHost
	} else {
		conf.DbHost = viper.GetString("db.host")
	}

	if dbPort := os.Getenv("DBPort"); len(dbPort) != 0 {
		conf.DbPort = dbPort
	} else {
		conf.DbPort = viper.GetString("db.port")
	}

	if dbLogin := os.Getenv("DBLogin"); len(dbLogin) != 0 {
		conf.DBLogin = dbLogin
	} else {
		conf.DBLogin = viper.GetString("db.login")
	}

	if dbPass := os.Getenv("DBPass"); len(dbPass) != 0 {
		conf.DBPass = dbPass
	} else {
		conf.DBPass = viper.GetString("db.pass")
	}

	if dbName := os.Getenv("DBName"); len(dbName) != 0 {
		conf.DbName = dbName
	} else {
		conf.DbName = viper.GetString("db.name")
	}

	return conf
}
