package Egolog

import (
	"easygo/conf"
	log "github.com/sirupsen/logrus"
	"os"
)

var Config Egoconf.Conf

var f *os.File

var LogFields = make(log.Fields)

func init() {
	Config = Egoconf.GetConfig()

	f = openLog(Config.LogPath, Config.LogName)

	logInit()

	Info("日志库完成初始化！")
}

func Info(msg string) {
	log.WithFields(LogFields).Info(msg)
}

func InfoData(msg string, data log.Fields) {
	for key, val := range LogFields {
		data[key] = val
	}
	log.WithFields(data).Info(msg)
}

func Fatal(msg string) {
	log.Fatal(msg)
}

func Error(msg string) {
	log.WithFields(LogFields).Error(msg)
}

func ErrorData(msg string, data log.Fields) {
	for key, val := range LogFields {
		data[key] = val
	}
	log.WithFields(data).Error(msg)
}

/**
 * 打开日志文件
 */
func openLog(logPath string, logName string) *os.File {

	if _, err := os.Stat(logPath); err != nil && os.IsNotExist(err) {
		if err := os.MkdirAll(logPath, os.ModePerm); err != nil {
			Error("日志目录创建失败!")
		}
	}
	f, err := os.OpenFile(logPath+logName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		Error("打开日志文件失败!")
	}

	return f
}

/**
 *日志基本配置
 */
func logInit() {
	log.SetOutput(f)

	if Egoconf.ENV == "production" {
		log.SetFormatter(&log.JSONFormatter{TimestampFormat: "2006-01-02 15:04:05"})
	} else {
		log.SetFormatter(&log.TextFormatter{DisableColors: true, FullTimestamp: true, TimestampFormat: "2006-01-02 15:04:05"})
	}
}
