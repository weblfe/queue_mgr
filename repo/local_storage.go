package repo

import (
	"errors"
	"fmt"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"github.com/weblfe/queue_mgr/utils"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

type localStorageRepository struct {
	defaultDB string
	locker    sync.RWMutex
	dbs       map[string]*leveldb.DB
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

func OpenDbWithEnv(name string) (*leveldb.DB, error) {
	var key = fmt.Sprintf("%s", "LOCAL_STORAGE_FILE")
	if name != "" {
		key = strings.ToUpper(fmt.Sprintf("%s_%s", name, key))
	}
	var file = utils.GetEnvVal(key)
	if file == "" {
		return nil, errors.New("miss db file")
	}
	var opts = &opt.Options{}
	if name == "" {
		opts.ReadOnly = utils.GetEnvBool("LOCAL_STORAGE_READ_ONLY")
	} else {
		opts.ReadOnly = utils.GetEnvBool(fmt.Sprintf("%s_%s", name, "LOCAL_STORAGE_READ_ONLY"))
	}
	return OpenDbWithFile(file, opts)
}

func OpenDbWithFile(file string, opts ...*opt.Options) (*leveldb.DB, error) {
	var dbFile, _ = filepath.Abs(file)
	if len(opts) <= 0 {
		opts = append(opts, &opt.Options{
			ReadOnly: utils.GetEnvBool("LOCAL_STORAGE_READ_ONLY"),
		})
	}
	if _, err := os.Stat(dbFile); err != nil {
		if opts[0].ReadOnly {
			return nil, err
		}
	}
	return leveldb.OpenFile(file, opts[0])
}

func CreateDbWithFile(file string, opts ...*opt.Options) (*leveldb.DB, error) {
	var dbFile, _ = filepath.Abs(file)
	if len(opts) <= 0 {
		opts = append(opts, &opt.Options{
			ReadOnly: utils.GetEnvBool("LOCAL_STORAGE_READ_ONLY"),
		})
	}
	if _, err := os.Stat(dbFile); err != nil {
		if !os.IsNotExist(err) {
			return nil, err
		}
		if err = os.MkdirAll(dbFile, 0755); err != nil {
			return nil, err
		}
	}
	return leveldb.OpenFile(file, opts[0])
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

func (repo *localStorageRepository) OpenDb(name string) error {
	if name == "" {
		return errors.New("miss db name")
	}
	repo.locker.Lock()
	defer repo.locker.Unlock()
	if v, ok := repo.dbs[name]; ok && v != nil {
		return nil
	}
	var db, err = OpenDbWithEnv(name)
	if err != nil {
		return err
	}
	if db == nil {
		return errors.New("open Db failed")
	}
	repo.dbs[name] = db
	return nil
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
