package logger

import (
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"path/filepath"
	"time"
)

type Info struct {
	LogLevel   string
	LogPath    string
	OutType    int // 0: 只输出到控制台 1: 只输出到日志
	TimeLayout string
}

// InitLogger 初始化日志配置
func InitLogger(info Info) {
	level, err := logrus.ParseLevel(info.LogLevel)
	if err != nil {
		panic(err)
	}
	logrus.SetLevel(level)
	logrus.SetReportCaller(true)
	var output io.Writer
	switch info.OutType {
	case 0:
		logrus.SetFormatter(&CustomStdoutFormatter{
			TimeLayout: info.TimeLayout,
		})
		output = os.Stdout
	case 1:
		logrus.SetFormatter(&CustomFileFormatter{
			TimeLayout: info.TimeLayout,
		})
		err = os.MkdirAll(info.LogPath, os.ModePerm)
		if err != nil {
			panic(err)
		}

		logFilePath := filepath.Join(info.LogPath, time.Now().Format("2006-01-02")+".log")
		file, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			logrus.Fatalf("Failed to open log file '%s': %v", logFilePath, err)
		}

		output = file
		go rotateLogDaily(file, info.LogPath)

	default:
		logrus.Fatalf("Invalid 'OutType' value: %d. Expected 0 or 1.", info.OutType)
	}
	logrus.SetOutput(output)
}

// rotateLogDaily handles daily log file rotation
func rotateLogDaily(currentFile *os.File, logPath string) {
	for {
		now := time.Now()
		nextDay := now.AddDate(0, 0, 1)
		nextDay = time.Date(nextDay.Year(), nextDay.Month(), nextDay.Day(), 0, 0, 0, 0, nextDay.Location())

		time.Sleep(time.Until(nextDay))

		currentFile.Close()
		logFilePath := filepath.Join(logPath, time.Now().Format("2006-01-02")+".log")
		newFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			logrus.Errorf("Failed to create new log file for the next day: %v", err)
			continue
		}
		logrus.SetOutput(newFile)
	}
}
