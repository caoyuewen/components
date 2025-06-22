package logger

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"path"
)

type CustomFileFormatter struct {
	TimeLayout string
}

func (f *CustomFileFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	timestamp := entry.Time.Format(f.TimeLayout)
	var file string
	var line int
	if entry.HasCaller() {
		file = path.Base(entry.Caller.File)
		line = entry.Caller.Line
	}

	var logFields string
	if len(entry.Data) > 0 {
		fields, err := json.Marshal(entry.Data)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal fields to JSON: %v", err)
		}
		logFields = string(fields)
	}
	level := fmtLevel(entry.Level, false)
	msg := getPrintLayout(level, timestamp, file, line, logFields, entry)
	return []byte(msg), nil
}
