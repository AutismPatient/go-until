package date

import (
	"fmt"
	"time"
)

const (
	defaultLayout  = "2006-01-02 15:04:05"
	sortLayout     = "2006-01-02"
	onlyLayout     = "01-02 15:04:05"
	nextLayout     = "01-02 15:04"
	monthDayLayout = "1-2"
)

// GetTimeText 将时间戳转换成string 长日期 2006-01-02 15:04:05
func GetTimeText(s int64) string {
	if s == 0 {
		return ""
	}
	return time.Unix(s, 0).Format(defaultLayout)
}

// GetTimeHistoryText 将时间戳转换成string 一天前、一个月前....
func GetTimeHistoryText(before int64) string {
	var now = time.Now().Unix()
	var offset = now - before
	switch {
	case offset >= 31536000:
		return GetSortTimeText(before)
	case offset >= 2592000:
		return fmt.Sprintf("%d月前", offset/2592000)
	case offset >= 86400:
		return fmt.Sprintf("%d天前", offset/86400)
	case offset < 86400:
		return fmt.Sprintf("%d小时前", offset/3600)
	}
	return "今天"
}

// GetDayTimeText 转换 -- 09-12 15:33:04
func GetDayTimeText(s int64) string {
	if s == 0 {
		return ""
	}
	return time.Unix(s, 0).Format(onlyLayout)
}

// GetDayTimeSortText 转换 -- 09-12 15:33
func GetDayTimeSortText(s int64) string {
	if s == 0 {
		return ""
	}
	return time.Unix(s, 0).Format(nextLayout)
}

// GetMonthDayText 转换 -- 9-12
func GetMonthDayText(s int64) string {
	if s == 0 {
		return ""
	}
	return time.Unix(s, 0).Format(monthDayLayout)
}

// GetSortTimeText 将时间戳转换成string 短日期
func GetSortTimeText(s int64) string {
	if s == 0 {
		return ""
	}
	return time.Unix(s, 0).Format(sortLayout)
}

// GetLastDateOfWeek 上一周 ，前2周，前3周 的周日|周结束
func GetLastDateOfWeek(d time.Time) time.Time {
	t := d.AddDate(0, 0, 6)
	return t
}

// GetFirstDateOfWeek 上一周 ，前2周，前3周 的周一|周开始
func GetFirstDateOfWeek(d time.Time) time.Time {
	offset := int(time.Monday - d.Weekday())
	if offset > 0 {
		offset = -7
	}
	t := d.AddDate(0, 0, offset)
	return t
}

// GetFirstDateOfMonth 获取上个月、上2个月、上3个月|月开始
func GetFirstDateOfMonth(d time.Time) time.Time {
	rr := d.AddDate(0, 0, -d.Day()+1)
	return rr
}
