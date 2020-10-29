package timer

import (
	"fmt"
	"github.com/robfig/cron/v3"
	string2 "go-until/string"
	"strings"
	"sync"
)

/*

 	cron表达式

	@yearly (or @annually) | Run once a year, midnight, Jan. 1st        | 0 0 1 1 *
	@monthly               | Run once a month, midnight, first of month | 0 0 1 * *
	@weekly                | Run once a week, midnight between Sat/Sun  | 0 0 * * 0
	@daily (or @midnight)  | Run once a day, midnight                   | 0 0 * * *
	@hourly



	@备注 参数：spec

	星号 （ * ）

	星号表示 cron 表达式将匹配字段的所有值;因此，cron 表达式将匹配字段的所有值。例如，在第 5 个字段（月）中使用星号将指示每个月。

	斜线 （ / ）

	斜线用于描述范围的增量。例如，第 1 个字段（分钟）中的 3-59/15 将指示小时的第 3 分钟以及此后每 15 分钟。表格"*\/..."等效于窗体"前一/..."，即字段最大可能范围的增量。"N/..."被接受为"N-MAX/..."，即从N开始，使用增量直到该特定范围的末尾。它不环绕。

	逗号 （ 、 ）

	逗号用于分隔列表的项。例如，在第 5 个字段（星期一的一天）中使用"MON、WED、FRI"意味着星期一、星期三和星期五。

	海芬 （ - ）

	连字符用于定义范围。例如，9-17 表示上午 9 点到下午 5 点（含）之间的每小时一次。

	问号 （ ？ ）

	可以使用问号代替"*"，以将月份或星期数留空。

	"@every 1h30m10s"表示在 1 小时 30 分钟 10 秒之后激活的时间表，然后在每个间隔之后激活

	* * * * * *
	Second | Minute | Hour | Dom | Month

	有关cron表达式：
	https://www.cnblogs.com/cangqinglang/p/10761567.html


	字段	 			允许值		允许的特殊字符

	秒（Seconds）	0~59的整数	, - * /    四个字符
	分（Minutes）	0~59的整数	, - * /    四个字符
	小时（Hours）	0~23的整数	, - * /    四个字符
	日期（DayofMonth）	1~31的整数（但是你需要考虑你月的天数）	,- * ? / L W C     八个字符
	月份（Month）	1~12的整数或者 JAN-DEC	, - * /    四个字符
	星期（DayofWeek）	1~7的整数或者 SUN-SAT （1=SUN）	, - * ? / L C #     八个字符
	年(可选，留空)（Year）	1970~2099	, - * /    四个字符


	常用表达式例子

　　（1）0 0 2 1 * ? *   表示在每月的1日的凌晨2点调整任务

　　（2）0 15 10 ? * MON-FRI   表示周一到周五每天上午10:15执行作业

　　（3）0 15 10 ? 6L 2002-2006   表示2002-2006年的每个月的最后一个星期五上午10:15执行作

　　（4）0 0 10,14,16 * * ?   每天上午10点，下午2点，4点

　　（5）0 0/30 9-17 * * ?   朝九晚五工作时间内每半小时

　　（6）0 0 12 ? * WED    表示每个星期三中午12点

　　（7）0 0 12 * * ?   每天中午12点触发

　　（8）0 15 10 ? * *    每天上午10:15触发

　　（9）0 15 10 * * ?     每天上午10:15触发

　　（10）0 15 10 * * ? *    每天上午10:15触发

　　（11）0 15 10 * * ? 2005    2005年的每天上午10:15触发

　　（12）0 * 14 * * ?     在每天下午2点到下午2:59期间的每1分钟触发

　　（13）0 0/5 14 * * ?    在每天下午2点到下午2:55期间的每5分钟触发

　　（14）0 0/5 14,18 * * ?     在每天下午2点到2:55期间和下午6点到6:55期间的每5分钟触发

　　（15）0 0-5 14 * * ?    在每天下午2点到下午2:05期间的每1分钟触发

　　（16）0 10,44 14 ? 3 WED    每年三月的星期三的下午2:10和2:44触发

　　（17）0 15 10 ? * MON-FRI    周一至周五的上午10:15触发

　　（18）0 15 10 15 * ?    每月15日上午10:15触发

　　（19）0 15 10 L * ?    每月最后一日的上午10:15触发

　　（20）0 15 10 ? * 6L    每月的最后一个星期五上午10:15触发

　　（21）0 15 10 ? * 6L 2002-2005   2002年至2005年的每月的最后一个星期五上午10:15触发

　　（22）0 15 10 ? * 6#3   每月的第三个星期五上午10:15触发

*/

type MyTimer struct {
	locker  sync.Mutex
	context *cron.Cron
	jobs    map[string]cron.EntryID
}

func NewTimer() *MyTimer {
	cn := cron.New(cron.WithLogger(cron.DefaultLogger), cron.WithSeconds())
	return &MyTimer{
		context: cn,
		jobs:    make(map[string]cron.EntryID),
	}
}
func (t *MyTimer) Start() {
	//if len(t.context.Entries()) > 0 {
	//
	//	return
	//}
	//log.Fatal("The timer must be started before the task is added")

	t.context.Start()
}
func (t *MyTimer) Stop() {
	t.context.Stop()
}

func (t *MyTimer) AddTask(spec string, f ...func()) (err error) {
	t.locker.Lock()
	defer t.locker.Unlock()

	id := string2.RandToken(6)

	if len(f) > 0 {
		for _, fc := range f {
			entityID, err := t.context.AddFunc(spec, fc)
			if err != nil {
				continue
			}
			t.jobs[id] = entityID
		}
	}

	if !(nil == err) {
		return err
	}

	return nil
}
func (t *MyTimer) AddJob(spec string, f ...cron.Job) (err error) {
	t.locker.Lock()
	defer t.locker.Unlock()

	id := string2.RandToken(6)

	if len(f) > 0 {
		for _, fc := range f {
			entityID, err := t.context.AddJob(spec, fc)
			if err != nil {
				continue
			}
			t.jobs[id] = entityID
		}
	}

	if !(nil == err) {
		return err
	}

	return nil
}

func SetHour(hour int) (spec string) {
	if hour >= 1 {
		spec = fmt.Sprintf("0/1 0/1 0/%d * * ?", hour)
	}
	return spec
}
func SetMinute(minute int) (spec string) {
	if minute > 0 && minute <= 60 {
		spec = fmt.Sprintf("0/1 0/%d * * * ?", minute)
	}
	return spec
}
func SetSecond(second int) (spec string) {
	if second > 0 && second <= 60 {
		spec = fmt.Sprintf("0/%d * * * * ?", second)
	}
	return spec
}

// 设置某个时刻下执行
func SetTimeOfDay(minute, hour int) (spec string) {
	spec = fmt.Sprintf("0 %d %d * * ?", minute, hour)
	return spec
}

// 在某个时间阶段性按间隔执行 参数限制格式如：3-6am,19-23pm
func SetSomeTimeOfGet(offset int, ams ...string) (spec string) {
	if len(ams) > 0 {
		spec = fmt.Sprintf("0 %d %s * * ?", offset, strings.Join(ams, ","))
	}
	return spec
}
