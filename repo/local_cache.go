package repo

import (
	"errors"
	"github.com/patrickmn/go-cache"
	"github.com/weblfe/queue_mgr/utils"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"
)

type (
	LocalCacheRepository struct {
		storage *cache.Cache
		options *Options
		remote  *RedisRepository
	}

	Options struct {
		DefaultExpiration time.Duration `json:"default_expiration" env:"local_cache_default_expiration"`
		CleanupInterval   time.Duration `json:"cleanup_interval" env:"local_cache_cleanup_interval"`
		SaveFile          string        `json:"save_file" env:"local_cache_save_file"`
		ReadOnly          bool          `json:"read_only" env:"local_cache_read_only"`
		Redis             string        `json:"redis" env:"local_cache_redis"`
	}
)

var (
	defaultOptions    *Options
	localCacheRepoIns *LocalCacheRepository
	locker            = sync.RWMutex{}
)

func CreateOptionsWithEnv() *Options {
	var opt = &Options{
		DefaultExpiration: 5 * time.Minute,
		CleanupInterval:   10 * time.Minute,
	}
	opt.ReadOnly = utils.GetEnvBool("LOCAL_CACHE_READ_ONLY")
	opt.SaveFile = utils.GetEnvVal("LOCAL_CACHE_STORAGE_FILE", "")
	opt.DefaultExpiration = utils.GetEnvDuration("LOCAL_CACHE_EXPIRATION", opt.DefaultExpiration)
	opt.CleanupInterval = utils.GetEnvDuration("LOCAL_CACHE_CLEANUP_INTERVAL", opt.CleanupInterval)
	return opt
}

func GetDefaultOptions() *Options {
	if defaultOptions == nil {
		defaultOptions = CreateOptionsWithEnv()
	}
	return defaultOptions
}

func newLocalCacheRepo(opts ...*Options) *LocalCacheRepository {
	var repo = new(LocalCacheRepository)
	opts = append(opts, GetDefaultOptions())
	repo.options = opts[0]
	return repo.load()
}

// GetLocalCacheRepo 获取本地缓存
func GetLocalCacheRepo() *LocalCacheRepository {
	if localCacheRepoIns == nil {
		locker.Lock()
		defer locker.Unlock()
		localCacheRepoIns = newLocalCacheRepo()
	}
	return localCacheRepoIns
}

func (repo *LocalCacheRepository) load() *LocalCacheRepository {
	if repo.storage == nil {
		repo.storage = cache.New(repo.options.DefaultExpiration, repo.options.CleanupInterval)
	}
	if repo.options.SaveFile != "" {
		var file, err = filepath.Abs(repo.options.SaveFile)
		if err != nil {
			repo.options.SaveFile = ""
		} else {
			repo.options.SaveFile = file
			if !repo.options.ReadOnly {
				_, err = os.Stat(file)
				if err != nil && os.IsNotExist(err) {
					base := filepath.Dir(file)
					_ = os.MkdirAll(base, os.ModePerm)
					if fs, err := os.Create(file); err == nil {
						_ = fs.Close()
					}
				}
			}
			if err = repo.storage.LoadFile(file); err != nil {
				if io.EOF != err {
					GetLogger("repo").Errorln(err)
				}
			}
			runtime.SetFinalizer(repo, (*LocalCacheRepository).Destroy)
		}
	}
	if repo.options.Redis != "" && repo.remote == nil {
		repo.remote = RedisDb(repo.options.Redis)
	}
	return repo
}

func (repo *LocalCacheRepository) GetCacheStorage() *cache.Cache {
	return repo.storage
}

func (repo *LocalCacheRepository) Set(key string, v interface{}, expire ...time.Duration) error {
	var store = repo.GetCacheStorage()
	if store == nil {
		return errors.New("open cache storage error")
	}
	if len(expire) <= 0 || expire[0] <= 0 {
		return store.Add(key, v, cache.NoExpiration)
	}
	return store.Add(key, v, expire[0])
}

func (repo *LocalCacheRepository) SetMust(key string, v interface{}, expire ...time.Duration) error {
	var store = repo.GetCacheStorage()
	if store == nil {
		return errors.New("open cache storage error")
	}
	_, ok := store.Get(key)
	if ok {
		if len(expire) <= 0 || expire[0] <= 0 {
			return store.Replace(key, v, cache.NoExpiration)
		}
		return store.Replace(key, v, expire[0])
	}
	if len(expire) <= 0 || expire[0] <= 0 {
		return store.Add(key, v, cache.NoExpiration)
	}
	return store.Add(key, v, expire[0])
}

func (repo *LocalCacheRepository) SetMustDefaultExpire(key string, v interface{}) error {
	var store = repo.GetCacheStorage()
	if store == nil {
		return errors.New("open cache storage error")
	}
	var _, ok = store.Get(key)
	if ok {
		return store.Replace(key, v, repo.options.DefaultExpiration)
	}
	return store.Add(key, v, repo.options.DefaultExpiration)
}

func (repo *LocalCacheRepository) Exists(key string) bool {
	var storage = repo.GetCacheStorage()
	if storage == nil {
		return false
	}
	_, ok := storage.Get(key)
	return ok
}

func (repo *LocalCacheRepository) Get(key string) (interface{}, bool) {
	var storage = repo.GetCacheStorage()
	if storage == nil {
		return nil, false
	}
	return storage.Get(key)
}

func (repo *LocalCacheRepository) remoteAdd(key string, v interface{}, expire ...time.Duration) {
	return
}

func (repo *LocalCacheRepository) remoteGet(key string) (interface{}, bool) {
	var storage = repo.getRemote()
	if storage == nil {
		return nil, false
	}
	return nil, false
}

func (repo *LocalCacheRepository) getRemote() *RedisRepository {
	return repo.remote
}

func (repo *LocalCacheRepository) Destroy() {
	defer runtime.SetFinalizer(repo, nil)
	repo.Sync()
	repo.storage = nil
	repo.options = nil
	repo.remote = nil
}

func (repo *LocalCacheRepository) Sync() {
	if repo.options.SaveFile != "" {
		var storage = repo.GetCacheStorage()
		if len(storage.Items()) <= 0 {
			return
		}
		if err := storage.SaveFile(repo.options.SaveFile); err != nil {
			GetLogger("repo").Errorln(err)
		}
	}
}
