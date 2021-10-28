package models

import (
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/weblfe/queue_mgr/repo"
	"time"
	"xorm.io/builder"
	"xorm.io/xorm"
	"xorm.io/xorm/names"
)

type (
	baseModel struct {
		table string
		cache string
		pk    string
	}
	M map[string]interface{}
)

var (
	logger *logrus.Logger
)

func (b *baseModel) TableName() string {
	var mapper = b.GetDb().GetTableMapper()
	if mapper != nil {
		return mapper.Obj2Table(b.table)
	}
	return b.table
}

func (b *baseModel) setTable(table string) {
	if b.table != "" || table == "" {
		return
	}
	b.table = table
}

func (b *baseModel) GetDb(name ...string) *xorm.Engine {
	name = append(name, "")
	var (
		key = name[0]
		db  = repo.GetDatabaseRepository().GetDb(key)
	)
	if db == nil {
		panic(errors.New("missing define connection: " + key))
	}
	engine, ok := db.(*xorm.Engine)
	if ok {
		return engine
	}
	if engine == nil {
		panic(errors.New("connection type error :" + fmt.Sprintf("%v", engine)))
	}
	return nil
}

func (b *baseModel) getTimezone() *time.Location {
	return b.GetDb().TZLocation
}

func (b *baseModel) Query() *xorm.Engine {
	return b.GetDb()
}

func (b *baseModel) save(data interface{}) (int64, error) {
	return b.Query().Insert(data)
}

func (b *baseModel) inserts(dataArr interface{}) (int64, error) {
	return b.Query().Insert(dataArr)
}

func (b *baseModel) getCache(name ...string) *repo.RedisRepository {
	name = append(name, b.cache)
	return repo.RedisDb(name[0])
}

func (b *baseModel) IsNotExits(err error) bool {
	if err == xorm.ErrNotExist {
		return true
	}
	if err == xorm.ErrObjectIsNil {
		return true
	}
	return false
}

func (b *baseModel) GetByID(id uint, model names.TableName) (interface{}, error) {
	if model == nil {
		return nil, errors.New("model params missing")
	}
	ok, err := b.Query().ID(id).Get(model)
	if ok && err == nil {
		return model, nil
	}
	return nil, err
}

func (b *baseModel) GetByIDS(id []uint, model names.TableName) (interface{}, error) {
	if model == nil {
		return nil, errors.New("model params missing")
	}
	ok, err := b.Query().Where(builder.In(b.getPkName(), id)).Get(model)
	if ok && err == nil {
		return model, nil
	}
	return nil, err
}

func (b *baseModel) Exists(cond builder.Cond, model names.TableName) bool {
	if model == nil {
		return false
	}
	var ok, err = b.Query().Table(model).Where(cond).Exist(model)
	if ok && err == nil {
		return ok
	}
	return false
}

// 获取主键名
func (b *baseModel) getPkName() string {
	if b.pk == "" {
		return "id"
	}
	return b.pk
}

// 设置主键名
func (b *baseModel) setPkName(pkName string) bool {
	if b.pk == "" && pkName != "" {
		b.pk = pkName
		return true
	}
	return false
}

func GetModelLogger() *logrus.Logger {
	if logger == nil {
		logger = repo.GetLogger("model")
	}
	return logger
}
