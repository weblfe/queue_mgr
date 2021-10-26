package entity

import (
	"errors"
	"github.com/weblfe/queue_mgr/utils"
	"sort"
	"time"
)

type (
	Crontab struct {
		At       int64          `json:"at"`
		TargetID string         `json:"targetID,omitempty"`
		Data     interface{}    `json:"data,omitempty"`
		Local    *time.Location `json:"local,omitempty"`
		callback func() error
	}

	CrontabItems []*Crontab
)

var (
	timezone     string
	timeOffset   int
	_local       *time.Location
	_PRCLocal, _ = time.LoadLocation("Asia/Shanghai")
)

func NewCrontab(at time.Time, id string, callback func() error) *Crontab {
	return &Crontab{At: parseTimestamp(at), callback: callback, TargetID: id, Local: getLocation()}
}

func parseTimestamp(at time.Time) int64 {
	return 0
}

func parseTime(t time.Time) *time.Time {
	if t.Unix() <= 0 {
		return nil
	}
	var (
		lock, err = time.Parse(utils.DateTimeLayout, t.Format(utils.DateTimeLayout))
	)
	if err != nil {
		return nil
	}
	var at = lock.In(getLocation())
	return &at
}

func getLocation() *time.Location {
	if _local != nil {
		return _local
	}
	// 时区
	if timezone == "" {
		timezone = utils.GetEnvVal("APP_TIME_ZONE", "")
	}
	// 时区偏远量
	if timeOffset <= 0 {
		timeOffset = utils.GetEnvInt("TIME_OFFSET", 0)
	}
	if _local == nil && timezone != "" {
		_local, _ = time.LoadLocation(timezone)
	}
	if _local == nil && timeOffset != 0 {
		_local = time.FixedZone("CST", timeOffset)
	}
	if _local == nil {
		_local = _PRCLocal
	}
	return _local
}

func (crontab *Crontab) SetLocal(local *time.Location) *Crontab {
	if crontab == nil || local == nil {
		return crontab
	}
	crontab.Local = local
	return crontab
}

func (crontab *Crontab) Execute() error {
	if crontab == nil {
		return errors.New("crontab nil")
	}
	if crontab.callback == nil {
		return nil
	}
	return crontab.callback()
}

func (crontab *Crontab) SetData(data interface{}) *Crontab {
	if crontab == nil {
		return crontab
	}
	if crontab.Data == nil {
		crontab.Data = data
	}
	return crontab
}

func (crontab *Crontab) ID() string {
	if crontab == nil {
		return ""
	}
	return crontab.TargetID
}

func (crontab *Crontab) Parse(now ...time.Time) bool {
	//var local = crontab.getLocal()
	//now = append(now, time.Now().In(local))
	now = append(now, time.Now())
	if crontab.At <= 0 {
		return true
	}
	var (
		t        = now[0]
		datetime = t.Format(utils.DateTimeLayout)
		clock, _ = time.Parse(utils.DateTimeLayout, datetime)
		ok       = clock.Unix() >= crontab.At
	)
	if ok {
		return true
	}
	return false
}

func (crontab *Crontab) getLocal() *time.Location {
	if crontab.Local == nil {
		crontab.Local = getLocation()
	}
	return crontab.Local
}

func (crontab *Crontab) Duration(now ...time.Time) time.Duration {
	var local = crontab.getLocal()
	now = append(now, time.Now().In(local))
	if crontab.At <= 0 {
		return 0
	}
	var t = now[0]
	return time.Duration(t.Unix() - crontab.At)
}

func (crontab *Crontab) DateTime() string {
	return time.Unix(crontab.At, 0).Format(utils.DateTimeLayout)
}

func (crontab *Crontab) Check() bool {
	if crontab.callback == nil {
		return false
	}
	return true
}

func (items CrontabItems) Len() int {
	return len(items)
}

func (items CrontabItems) Less(i, j int) bool {
	if items[i].At <= items[j].At {
		return true
	}
	return false
}

func (items CrontabItems) Swap(i, j int) {
	items[i], items[j] = items[j], items[i]
}

func (items *CrontabItems) Append(it *Crontab) *CrontabItems {
	*items = append(*items, it)
	return items
}

func (items *CrontabItems) Pop() *Crontab {
	var (
		it = (*items)[0]
	)
	*items = (*items)[1:]
	return it
}

func (items *CrontabItems) Insert(it *Crontab) *CrontabItems {
	items.Append(it)
	sort.Sort(items)
	return items
}

func (items CrontabItems) List() []time.Time {
	var (
		now = time.Now()
		arr []time.Time
	)
	for _, v := range items {
		if v.At > 0 {
			arr = append(arr, time.Unix(v.At, 0))
		} else {
			arr = append(arr, now)
		}
	}
	return arr
}
