package util

import "time"

func InSameDay(t1, t2 time.Time, loc *time.Location) bool {
	t1 = t1.In(loc)
	t2 = t2.In(loc)
	return t1.Year() == t2.Year() && t1.Month() == t2.Month() && t1.Day() == t2.Day()
}

func InSameMonth(t1, t2 time.Time, loc *time.Location) bool {
	t1 = t1.In(loc)
	t2 = t2.In(loc)
	return t1.Year() == t2.Year() && t1.Month() == t2.Month()
}

// GetMondayZeroTime 获取本周一 0 点时间
func GetMondayZeroTime() time.Time {
	now := time.Now()
	// 获取本地时区
	//_, offset := now.Zone()
	// 获取今天是星期几，Go 的 time.Weekday() 返回 0 (Sunday) 到 6 (Saturday)
	weekday := int(now.Weekday())

	// 让周一为一周的第一天。如果是星期日(0)，则往前推6天，否则往前推 weekday-1 天
	if weekday == 0 {
		weekday = 7 // 将 Sunday 视为第7天
	}

	// 计算本周一的时间
	startOfWeek := now.AddDate(0, 0, -(weekday - 1))
	// 设置为当天的 0 点
	startOfWeek = time.Date(startOfWeek.Year(), startOfWeek.Month(), startOfWeek.Day(), 0, 0, 0, 0, time.Local)
	// 考虑时区的偏移，减去 offset 秒
	//startOfWeek = startOfWeek.Add(time.Duration(-offset) * time.Second)
	return startOfWeek
}

func GetNextMondayZeroTime(t time.Time) time.Time {
	// 获取当前是星期几，Go 的 time.Weekday() 返回 0 (Sunday) 到 6 (Saturday)
	weekday := int(t.Weekday())

	// 如果今天是周一且已过了 0 点，则下一次的周一为一周后的周一
	daysUntilNextMonday := (7 - weekday) % 7
	if daysUntilNextMonday == 0 && t.Hour() >= 0 {
		daysUntilNextMonday = 7
	}

	// 计算下一个周一的日期
	nextMonday := t.AddDate(0, 0, daysUntilNextMonday)

	// 返回下一个周一 0 点的时间
	return time.Date(nextMonday.Year(), nextMonday.Month(), nextMonday.Day(), 0, 0, 0, 0, t.Location())
}
