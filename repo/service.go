package repo

import (
	"errors"
	log "github.com/sirupsen/logrus"
	"github.com/weblfe/queue_mgr/facede"
	"sync"
)

type serverRegisterRepository struct {
	cache map[string]facede.Service
	lock  sync.RWMutex
}

var (
	serverRegister    = newServerRegisterRepository()
	ErrorServNotFound = errors.New("service not found")
)

func newServerRegisterRepository() *serverRegisterRepository {
	var repo = new(serverRegisterRepository)
	repo.lock = sync.RWMutex{}
	repo.cache = make(map[string]facede.Service)
	return repo
}

// GetServerRegisterRepository 获取服务注册库
func GetServerRegisterRepository() *serverRegisterRepository {
	return serverRegister
}

func (repo *serverRegisterRepository) Get(key string) (facede.Service, error) {
	repo.lock.Lock()
	defer repo.lock.Unlock()
	if v, ok := repo.cache[key]; ok && v != nil {
		return v, nil
	}
	return nil, ErrorServNotFound
}

func (repo *serverRegisterRepository) Exists(key string) bool {
	repo.lock.RLocker().Lock()
	defer repo.lock.RLocker().Unlock()
	if _, ok := repo.cache[key]; ok {
		return true
	}
	return false
}

func (repo *serverRegisterRepository) Register(key string, service interface{}, force ...bool) bool {
	if key == "" || service == nil {
		return false
	}
	s, ok := service.(facede.Service)
	if !ok {
		return false
	}
	force = append(force, false)
	repo.lock.RLocker().Lock()
	defer repo.lock.RLocker().Unlock()
	v, ok := repo.cache[key]
	if ok {
		if v == nil || force[0] {
			repo.cache[key] = s
			return true
		}
		return false
	}
	repo.cache[key] = s
	return true
}

func (repo *serverRegisterRepository) Load(service facede.Service, force ...bool) bool {
	force = append(force, false)
	repo.lock.RLocker().Lock()
	defer repo.lock.RLocker().Unlock()
	key := service.ServiceID()
	v, ok := repo.cache[key]
	if ok {
		if v != nil && !force[0] {
			return false
		}
	}
	repo.cache[key] = service
	log.Infoln("注册服务:", key, "成功")
	return true
}


func (repo *serverRegisterRepository) GetUserService() facede.UserService {
	repo.lock.RLocker().Lock()
	defer repo.lock.RLocker().Unlock()
	var key = "userService"
	v, ok := repo.cache[key]
	if !ok {
		return nil
	}
	serv, ok := v.(facede.UserService)
	if ok {
		return serv
	}
	return nil
}
