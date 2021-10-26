package repo

import (
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/weblfe/queue_mgr/utils"
	"runtime"
	"sync"
)

type localStorageRepository struct {
	locker    sync.RWMutex
	dbs       map[string]*leveldb.DB
	defaultDB string
}

var (
	localStorage = newLocalStorageRepo()
)

func newLocalStorageRepo() *localStorageRepository {
	var storage = &localStorageRepository{}

	return storage.init()
}

func GetLocalStorageRepo() *localStorageRepository {
	return localStorage
}

func (repo *localStorageRepository) init() *localStorageRepository {
	repo.locker = sync.RWMutex{}
	repo.dbs = make(map[string]*leveldb.DB)
	repo.defaultDB = utils.GetEnvVal("LOCAL_STORAGE_DEFAULT", "default")
	// 对象回收时
	runtime.SetFinalizer(repo, (*localStorageRepository).destroy)
	return repo
}

func (repo *localStorageRepository) RegisterDb(name string, db *leveldb.DB) bool {
	if db == nil || name == "" {
		return false
	}
	repo.locker.Lock()
	defer repo.locker.Unlock()
	if v, ok := repo.dbs[name]; ok && v != nil {
		return false
	}
	repo.dbs[name] = db
	return true
}

func (repo *localStorageRepository) Remove(name string) bool {
	if name == "" {
		return false
	}
	repo.locker.Lock()
	defer repo.locker.Unlock()
	v, ok := repo.dbs[name]
	if !ok {
		return false
	}
	delete(repo.dbs, name)
	if err := v.Close(); err != nil {
		GetLogger("local_storage").WithField("error", err.Error()).Errorln("local storage db:" + name + ",close error")
	}
	return true
}

func (repo *localStorageRepository) getDbs() []string {
	var dbs []string
	for k := range repo.dbs {
		dbs = append(dbs, k)
	}
	return dbs
}

func (repo *localStorageRepository) destroy() error {
	repo.locker.Lock()
	defer repo.locker.Unlock()
	runtime.SetFinalizer(repo, nil)
	var (
		errs []error
		dbs  = repo.getDbs()
	)
	if len(dbs) <= 0 {
		return nil
	}
	for _, key := range dbs {
		db, ok := repo.dbs[key]
		if !ok {
			continue
		}
		delete(repo.dbs, key)
		if err := db.Close(); err != nil {
			GetLogger("local_storage").WithField("error", err.Error()).Errorln("local storage db:" + key + ",close error")
			errs = append(errs, err)
		}
	}
	if len(errs) > 0 {
		return errs[0]
	}
	return nil
}

func (repo *localStorageRepository) GetStorage(name ...string) (*leveldb.DB, bool) {
	name = append(name, repo.defaultDB)
	repo.locker.Lock()
	defer repo.locker.Unlock()
	var key = name[0]
	if db, ok := repo.dbs[key]; ok {
		return db, true
	}
	return nil, false
}
