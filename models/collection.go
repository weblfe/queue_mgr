package models

import (
	"reflect"
	"sync"
	"xorm.io/xorm"
)

type collectionImpl struct {
	iter        func(v interface{})
	item        interface{} // 单个对象
	rows        *xorm.Rows
	constructor sync.Once
	errs        []error
	cache       []interface{}
	size        int
}

func NewCollection(rows *xorm.Rows, item interface{}) *collectionImpl {
	return &collectionImpl{
		constructor: sync.Once{},
		rows:        rows,
		item:        item,
	}
}

func (impl *collectionImpl) CreateIter(iter func(v interface{})) *collectionImpl {
	if impl.iter != nil {
		return impl
	}
	impl.iter = iter
	return impl
}

func (impl *collectionImpl) Errors() []error {
	return impl.errs
}

func (impl *collectionImpl) LastErr() error {
	if len(impl.errs) > 0 {
		return impl.errs[len(impl.errs)-1]
	}
	return nil
}

func (impl *collectionImpl) Parse() error {
	impl.constructor.Do(impl.createIterFunc())
	if len(impl.errs) > 0 {
		return impl.errs[len(impl.errs)-1]
	}
	return nil
}

func (impl *collectionImpl) createIterFunc() func() {
	var (
		rows = impl.rows
		item = impl.item
		iter = impl.iter
	)
	return func() {
		defer impl.close()
		for rows.Next() {
			impl.size++
			if err := rows.Scan(item); err != nil {
				GetModelLogger().Errorln("error:", err)
				impl.errs = append(impl.errs, err)
				continue
			}
			iter(item)
			// impl.append(item)
		}
	}
}

func (impl *collectionImpl) Count() int {
	return impl.size
}

func (impl *collectionImpl) append(v interface{}) {
	var value = reflect.ValueOf(v)
	if value.CanSet() {
		impl.cache = append(impl.cache, value.Elem().Interface())
	}else {
		impl.cache = append(impl.cache, value.Interface())
	}
}

func (impl *collectionImpl) close() {
	if impl.rows == nil {
		return
	}
	if err := impl.rows.Close(); err != nil {
		GetModelLogger().Errorln("collection close.error:", err)
		impl.errs = append(impl.errs, err)
		impl.rows = nil
	}
}

func (impl *collectionImpl) Reset() {
	impl.size = 0
	impl.rows = nil
	impl.item = nil
	impl.iter = nil
	impl.errs = nil
	impl.cache = nil
	impl.constructor = sync.Once{}
}
