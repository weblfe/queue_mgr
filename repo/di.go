package repo

import (
	"errors"
	"runtime"
	"sort"
	"sync"
)

type containerRepository struct {
	diTree map[string]interface{}
	locker sync.RWMutex
}

var (
	di               = newContainerRepo()
	ErrorNotFound    = errors.New("not found")
	ErrorNil         = errors.New("value nil")
	ErrorType        = errors.New("type error")
	ErrorResolverNil = errors.New("nil func resolver")
)

func newContainerRepo() *containerRepository {
	var container = &containerRepository{
		locker: sync.RWMutex{},
		diTree: make(map[string]interface{}),
	}
	container.init()
	return container
}

func (di *containerRepository) init() {
	runtime.SetFinalizer(di, (*containerRepository).destroy)
}

func GetContainerRepo() *containerRepository {
	return di
}

func (di *containerRepository) Get(key string, def ...interface{}) (interface{}, bool) {
	di.locker.RLocker().Lock()
	defer di.locker.RLocker().Unlock()
	def = append(def, nil)
	var v, ok = di.diTree[key]
	if !ok {
		return def[0], false
	}
	return v, ok
}

func (di *containerRepository) Exists(key string) bool {
	di.locker.RLocker().Lock()
	defer di.locker.RLocker().Unlock()
	var _, ok = di.diTree[key]
	return ok
}

func (di *containerRepository) Store(key string, v interface{}) bool {
	di.locker.RLocker().Lock()
	defer di.locker.RLocker().Unlock()
	if _, ok := di.diTree[key]; !ok {
		di.diTree[key] = v
		return true
	}
	return false
}

func (di *containerRepository) Remove(key string) bool {
	di.locker.RLocker().Lock()
	defer di.locker.RLocker().Unlock()
	if _, ok := di.diTree[key]; !ok {
		return false
	}
	delete(di.diTree, key)
	return true
}

func (di *containerRepository) Len() int {
	return len(di.diTree)
}

func (di *containerRepository) Keys() []string {
	var keys []string
	for k := range di.diTree {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func (di *containerRepository) Add(key string, v interface{}, force ...bool) bool {
	di.locker.RLocker().Lock()
	defer di.locker.RLocker().Unlock()
	force = append(force, false)
	if _, ok := di.diTree[key]; !ok || force[0] {
		di.diTree[key] = v
		return true
	}
	return false
}

func (di *containerRepository) Cache(key string, v interface{}) bool {
	if v == nil {
		return false
	}
	return di.Store(key, v)
}

func (di *containerRepository) Register(key string, v interface{}) *containerRepository {
	di.Store(key, v)
	return di
}

func (di *containerRepository) Resolve(key string, resolver func(v interface{}) error) error {
	var v, ok = di.Get(key)
	if !ok {
		return ErrorNotFound
	}
	if resolver == nil {
		return ErrorResolverNil
	}
	return resolver(v)
}

func (di *containerRepository) destroy() {
	runtime.SetFinalizer(di, nil)
	di.locker.Lock()
	defer di.locker.Unlock()
	for _, key := range di.Keys() {
		delete(di.diTree, key)
	}
	di.diTree = nil
}
