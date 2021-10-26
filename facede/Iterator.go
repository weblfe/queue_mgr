package facede

import "time"

// Iterator 迭代器
type Iterator interface {
	// Next 是否可以继续迭代
	Next() bool
	// Key 获取当前游标位置键名
	Key() string
	// Value 获取当前游标位置数据
	Value() interface{}
	// Item 获取当前key 和 值
	Item() (key string, value interface{})
	// Offset 获取对应位置数据
	Offset(i int) interface{}
	// Len 迭代容量
	Len() int
	// Cursor 移动游标
	Cursor() int
	// Reset 重置迭代
	Reset() bool
	// Foreach 遍历
	Foreach(each func(key string, value interface{}) bool)
	// ForeachArr 带数组下标 方式遍历
	ForeachArr(each func(i int, key string, value interface{}) bool)
}

type ValueConverter interface {
	String() string
	Interface() interface{}
	Int() int
	Duration() time.Duration
	Time() (time.Time, bool)
	Float() float64
	IsNull() bool
	IsZero() bool
	Empty() bool
}
