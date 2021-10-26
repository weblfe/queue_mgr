package entity

import (
	"fmt"
	"sort"
	"testing"
	"time"
)

func TestCrontabItems_Swap(t *testing.T) {
	var (
		now, _       = time.Parse(`2006-01-02 15:04:05`, "2021-09-27 11:13:43")
		crontabItems = CrontabItems{
			&Crontab{At: now.Unix()},
			&Crontab{At: now.Add(time.Minute).Unix()},
			&Crontab{At: now.Add(-time.Minute).Unix()},
			&Crontab{At: now.Add(-time.Hour).Unix()},
			&Crontab{At: now.Add(time.Hour).Unix()},
			&Crontab{At: now.Add(-2 * time.Hour).Unix()},
			&Crontab{At: now.Add(-3 * time.Hour).Unix()},
			&Crontab{At: now.Add(-3 * time.Hour).Unix()},
			&Crontab{At: now.Add(-3 * time.Hour).Unix()},
			&Crontab{At: now.Add(-3 * time.Hour).Unix()},
			&Crontab{At: now.Add(-3 * time.Hour).Unix()},
			&Crontab{At: now.Add(-3 * time.Hour).Unix()},
			&Crontab{At: now.Add(-3 * time.Hour).Unix()},
		}
	)
	sort.Sort(crontabItems)
	if crontabItems[0].At > crontabItems[len(crontabItems)-1].At {
		t.Error("CrontabItems Sort Error")
	}
}

func newTime(now *time.Time, duration time.Duration) *time.Time {
	var t = now.Add(duration)
	return &t
}

func TestCrontabItems_Insert(t *testing.T) {
	var (
		now, _       = time.Parse(`2006-01-02 15:04:05`, "2021-09-27 11:13:43")
		crontabItems = CrontabItems{}
	)
	crontabItems.Insert(&Crontab{At: now.Unix()}).
		Insert(&Crontab{At: now.Add(time.Minute).Unix()}).
		Insert(&Crontab{At: now.Add(-time.Minute).Unix()}).
		Insert(&Crontab{At: now.Add(-time.Hour).Unix()}).
		Insert(&Crontab{At: now.Add(1000 * time.Hour).Unix()}).
		Insert(&Crontab{At: now.Add(-2 * time.Hour).Unix()}).
		Insert(&Crontab{At: now.Add(-3 * time.Hour).Unix()}).
		Insert(&Crontab{At: now.Add(-3 * time.Hour).Unix()}).
		Insert(&Crontab{At: now.Add(-3 * time.Hour).Unix()}).
		Insert(&Crontab{At: now.Add(-3 * time.Hour).Unix()}).
		Insert(&Crontab{At: now.Add(-3 * time.Hour).Unix()}).
		Insert(&Crontab{At: now.Add(-3 * time.Hour).Unix()}).
		Insert(&Crontab{At: now.Add(-3 * time.Hour).Unix()})
	if crontabItems[0].At > crontabItems[len(crontabItems)-1].At {
		t.Error("CrontabItems Sort Error")
	}
	fmt.Println(crontabItems.Len())
	fmt.Println(crontabItems.List())
}

func TestCrontabItems_Pop(t *testing.T) {
	var (
		now, _       = time.Parse(`2006-01-02 15:04:05`, "2021-09-27 11:13:43")
		crontabItems = CrontabItems{}
	)
	crontabItems.Insert(&Crontab{At: now.Unix()}).
		Insert(&Crontab{At: now.Add(time.Minute).Unix()}).
		Insert(&Crontab{At: now.Add(-time.Minute).Unix()}).
		Insert(&Crontab{At: now.Add(-time.Hour).Unix()}).
		Insert(&Crontab{At: now.Add(1000 * time.Hour).Unix()}).
		Insert(&Crontab{At: now.Add(-2 * time.Hour).Unix()}).
		Insert(&Crontab{At: now.Add(-3 * time.Hour).Unix()}).
		Insert(&Crontab{At: now.Add(-3 * time.Hour).Unix()}).
		Insert(&Crontab{At: now.Add(-3 * time.Hour).Unix()}).
		Insert(&Crontab{At: now.Add(-3 * time.Hour).Unix()}).
		Insert(&Crontab{At: now.Add(-3 * time.Hour).Unix()}).
		Insert(&Crontab{At: now.Add(-3 * time.Hour).Unix()}).
		Insert(&Crontab{At: now.Add(-3 * time.Hour).Unix()})

	if crontabItems[0].At > crontabItems[len(crontabItems)-1].At {
		t.Error("CrontabItems Sort Error")
	}
	for crontabItems.Len() > 0 {
		it := crontabItems.Pop()
		if it == nil {
			t.Error("pop Error")
		}else {
			fmt.Println(it.DateTime())
		}
	}
}

func TestCrontab_Parse(t *testing.T) {
	var (
		// now, _       = time.Parse(`2006-01-02 15:04:05`, "2021-09-27 11:13:43")
		now          = time.Now()
		crontabItems = CrontabItems{}
	)
	crontabItems.Insert(&Crontab{At: now.Unix()}).
		Insert(&Crontab{At: now.Add(time.Minute).Unix()}).
		Insert(&Crontab{At: now.Add(-time.Minute).Unix()}).
		Insert(&Crontab{At: now.Add(-time.Hour).Unix()}).
		Insert(&Crontab{At: now.Add(1000 * time.Hour).Unix()}).
		Insert(&Crontab{At: now.Add(-2 * time.Hour).Unix()}).
		Insert(&Crontab{At: now.Add(-3 * time.Hour).Unix()}).
		Insert(&Crontab{At: now.Add(-3 * time.Hour).Unix()}).
		Insert(&Crontab{At: now.Add(-3 * time.Hour).Unix()}).
		Insert(&Crontab{At: now.Add(-3 * time.Hour).Unix()}).
		Insert(&Crontab{At: now.Add(-3 * time.Hour).Unix()}).
		Insert(&Crontab{At: now.Add(-3 * time.Hour).Unix()}).
		Insert(&Crontab{At: now.Add(-3 * time.Hour).Unix()})

	if crontabItems[0].At > crontabItems[len(crontabItems)-1].At {
		t.Error("CrontabItems Sort Error")
	}
	for crontabItems.Len() > 0 {
		it := crontabItems.Pop()
		if it == nil {
			t.Error("pop Error")
		}
		if it.Parse() {
			fmt.Println("时间: ", it.DateTime())
		}
	}
}
