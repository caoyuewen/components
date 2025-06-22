package util

import (
	"time"
)

// 一年中第几周，[0..53]
// 当firstDayOfWeek = time.Sunday，相当于date +%U
// 当firstDayOfWeek = time.Monday，相当于date +%W
func GetWeek(t time.Time, firstDayOfWeek time.Weekday) int {
	yearDay := t.YearDay()
	yearFirstDay := t.AddDate(0, 0, 1-yearDay)

	firstWeekDays := (int(firstDayOfWeek) - int(yearFirstDay.Weekday()) + 7) % 7
	week := (yearDay + 6 - firstWeekDays) / 7
	return week
}

// 是否为今天
func IsTodayTimestamp(timestamp int64, tz *time.Location) bool {
	now := time.Now().In(tz)
	t := time.Unix(timestamp, 0).In(tz)
	return now.Year() == t.Year() && now.Month() == t.Month() && now.Day() == t.Day()
}

// 判断两个时间戳相差多少天
func IntervalDaysTimestamp(ts1, ts2 int64, tz *time.Location) int {
	t1 := time.Unix(ts1, 0).In(tz)
	t2 := time.Unix(ts2, 0).In(tz)
	t1 = time.Date(t1.Year(), t1.Month(), t1.Day(), 0, 0, 0, 0, tz)
	t2 = time.Date(t2.Year(), t2.Month(), t2.Day(), 0, 0, 0, 0, tz)
	return int(t2.Sub(t1) / (24 * time.Hour))
}

func GetTodayTimestamp(loc *time.Location) int64 {
	// 获取当前时间
	currentTime := time.Now().In(loc)
	// 获取今日零点的时间
	today := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 0, 0, 0, 0, currentTime.Location())
	// 获取今日零点的时间戳
	todayTimestamp := today.Unix()
	return todayTimestamp
}
