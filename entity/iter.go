package entity

import (
	"sync"
)

type (
	objectIterImpl struct {
		keys       Keys
		values     Values
		safeLocker sync.RWMutex
		cursor     uint
		size       int
		capacity   int
	}

	Keys   []string
	Values []interface{}
)

func NewIterImpl() *objectIterImpl {
	var iter = new(objectIterImpl)
	iter.safeLocker = sync.RWMutex{}
	return iter
}

const (
	IterEoL = -1
)

func (iter *objectIterImpl) SetValues(values []interface{}) *objectIterImpl {
	if iter == nil {
		return nil
	}
	if len(values) <= 0 {
		return iter
	}
	iter.safeLocker.Lock()
	defer iter.safeLocker.Unlock()
	if len(iter.values) <= 0 {
		iter.values = values
	} else {
		iter.values = append(iter.values, values...)
	}
	iter.cap()
	return iter
}

func (iter *objectIterImpl) SetKeys(keys []string) *objectIterImpl {
	if iter == nil {
		return nil
	}
	if len(keys) <= 0 {
		return iter
	}
	iter.safeLocker.Lock()
	defer iter.safeLocker.Unlock()
	if len(iter.keys) <= 0 {
		iter.keys = keys
	} else {
		iter.keys = append(iter.keys, keys...)
	}
	iter.resize()
	return iter
}

func (iter *objectIterImpl) Len() int {
	if iter != nil && iter.size == 0 {
		iter.resize()
	}
	return iter.size
}

func (iter *objectIterImpl) Next() bool {
	if iter != nil && iter.cursor < uint(len(iter.keys)) {
		return true
	}
	return false
}

func (iter *objectIterImpl) Key() string {
	var cursor = iter.GetCursor()
	if iter.Len() > cursor {
		return iter.keys[cursor]
	}
	return ""
}

func (iter *objectIterImpl) Offset(i int) interface{} {
	if iter.Capacity() > i {
		return NewValue(iter.values[i])
	}
	return nil
}

func (iter *objectIterImpl) OffsetKey(i int) string {
	if iter.Len() > i {
		return iter.keys[i]
	}
	return ""
}

func (iter *objectIterImpl) Exists(k interface{}) bool {
	if k == nil {
		return false
	}
	switch k.(type) {
	case string:
		return iter.existKey(k.(string))
	case int:
		return iter.Len() > k.(int)
	case uint:
		return iter.Len() > int(k.(uint))
	case int32:
		return iter.Len() > int(k.(int32))
	case int64:
		return iter.Len() > int(k.(int64))
	}
	return false
}

func (iter *objectIterImpl) existKey(key string) bool {
	iter.safeLocker.Lock()
	defer iter.safeLocker.Unlock()
	if len(iter.keys) <= 0 {
		return false
	}
	for _, k := range iter.keys {
		if k == key {
			return true
		}
	}
	return false
}

func (iter *objectIterImpl) Value() interface{} {
	var cursor = iter.GetCursor()
	if iter.Capacity() > cursor {
		return NewValue(iter.values[cursor])
	}
	return nil
}

func (iter *objectIterImpl) GetCursor() int {
	return int(iter.cursor)
}

func (iter *objectIterImpl) Capacity() int {
	if iter != nil && iter.capacity <= 0 {
		iter.cap()
	}
	return iter.capacity
}

func (iter *objectIterImpl) Cursor() int {
	iter.safeLocker.Lock()
	defer iter.safeLocker.Unlock()
	if iter.Next() {
		iter.cursor++
	} else {
		return IterEoL
	}
	return int(iter.cursor)
}

func (iter *objectIterImpl) Foreach(each func(key string, value interface{}) bool) {
	iter.safeLocker.Lock()
	defer iter.safeLocker.Unlock()
	if iter.Len() <= 0 || each == nil {
		return
	}
	for i, k := range iter.keys {
		if !each(k, iter.Offset(i)) {
			break
		}
	}
}

func (iter *objectIterImpl) ForeachArr(each func(i int, key string, value interface{}) bool) {
	iter.safeLocker.Lock()
	defer iter.safeLocker.Unlock()
	if iter.Len() <= 0 || each == nil {
		return
	}
	for i, v := range iter.values {
		if !each(i, iter.OffsetKey(i), v) {
			break
		}
	}
}

func (iter *objectIterImpl) Reset() bool {
	iter.safeLocker.Lock()
	defer iter.safeLocker.Unlock()
	iter.cursor = 0
	iter.size = len(iter.keys)
	return true
}

func (iter *objectIterImpl) Item() (key string, value interface{}) {
	iter.safeLocker.Lock()
	defer iter.safeLocker.Unlock()
	return iter.Key(), iter.Value()
}

func (iter *objectIterImpl) resize() {
	iter.size = len(iter.keys)
}

func (iter *objectIterImpl) cap() {
	iter.capacity = len(iter.values)
}
