package logger

import (
	"os"
	"path"

	"github.com/DarkSoul94/money-processing-service/config"
	"github.com/sirupsen/logrus"
)

// InitLogger ...
func InitLogger(conf config.Config) {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	os.Mkdir(conf.LogDir, 0755)
	filepath := path.Join(conf.LogDir, conf.LogFile)
	file, err := os.OpenFile(filepath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		panic(err)
	}
	logrus.SetOutput(file)
}

// LogError log it
func LogError(action, file, data string, err error) {
	logrus.WithFields(
		logrus.Fields{
			"action": action,
			"file":   file,
			"data":   data,
		},
	).Error(err)
}
