package logger

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"strings"
)

func getPrintLayout(level, timestamp, file string, line int, logFields string, entry *logrus.Entry) string {
	// 固定文件名占用20个字符，行号占用4个字符
	if len(file) > 12 {
		file = fmt.Sprintf("%s_%s", file[:6], file[len(file)-5:])
	}
	formattedFile := fmt.Sprintf("%-12s", file)
	formattedLine := fmt.Sprintf("%4d", line)

	return fmt.Sprintf("%s [%s]%s %s %s %s\n", timestamp, formattedFile, formattedLine, level, entry.Message, logFields)
}

func fmtLevel(l logrus.Level, errHLight bool) string {
	level := l.String()
	switch level {
	case "warning":
		level = "warn"
	case "error":
		level = "erro"
	case "debug":
		level = "dbug"
	case "trace":
		level = "trac"
	}
	level = strings.ToUpper(level)
	if l <= logrus.ErrorLevel && errHLight {
		bold := "\033[1m"
		reset := "\033[0m"
		level = fmt.Sprintf("%s[%s]%s", bold, level, reset)
	} else {
		level = fmt.Sprintf("[%s]", level)
	}
	return level
}
