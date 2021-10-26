package repo

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
	"github.com/weblfe/queue_mgr/config"
	"github.com/weblfe/queue_mgr/utils"
	"os"
	"time"
	"xorm.io/xorm"
	"xorm.io/xorm/log"
	"xorm.io/xorm/names"
)

type databaseRepository struct {
	connection config.DatabaseKv
	database   map[string]xorm.EngineInterface
}

var (
	defaultDatabaseRepo = newDatabaseRepository()
)

func GetDatabaseRepository() *databaseRepository {
	return defaultDatabaseRepo
}

func newDatabaseRepository() *databaseRepository {
	var repo = new(databaseRepository)
	repo.connection = nil
	repo.database = make(map[string]xorm.EngineInterface)
	return repo
}

func (repo *databaseRepository) InitConnection(conn config.DatabaseKv) *databaseRepository {
	if repo.connection == nil {
		repo.connection = conn
	}
	return repo
}

func (repo *databaseRepository) GetConn(name string) config.Database {
	if repo.connection == nil {
		return config.GetAppConfig().GetDatabase(name)
	}
	var db, _ = repo.connection.Get(name)
	return db
}

func (repo *databaseRepository) GetDb(name ...string) xorm.EngineInterface {
	name = append(name, repo.getDefaultConnKey())
	var key = name[0]
	if conn, ok := repo.database[key]; ok && conn != nil {
		return conn
	}
	db := repo.GetConn(key)
	conn, err := xorm.NewEngine(db.DbDriver, db.GetConnUrl())
	// sql
	conn.SetLogger(log.NewSimpleLogger(GetLogger("sql").Out))
	// sql debug
	if utils.GetEnvBool("DB_SQL_DEBUG") {
		conn.Logger().ShowSQL(true)
	}
	// Database Timezone
	if local, err2 := time.LoadLocation(os.Getenv("TZ")); err2 == nil && local != nil {
		conn.DatabaseTZ = local
	}
	// App Timezone
	if local, err2 := time.LoadLocation(os.Getenv("APP_TIMEZONE")); err2 == nil && local != nil {
		conn.TZLocation = local
	}
	// 表前缀
	if db.DbPrefix != "" {
		conn.SetTableMapper(names.NewPrefixMapper(names.SameMapper{}, db.DbPrefix))
	}
	if err != nil {
		logrus.Infoln("Db Conn Error:", err)
		panic(err)
	}
	repo.Add(key, conn)
	return conn
}

func (repo *databaseRepository) getDefaultConnKey() string {
	var name = os.Getenv("DB_DEFAULT_CONN_NAME")
	if name == "" {
		return repo.connection.GetDefaultKey()
	}
	return name
}

func (repo *databaseRepository) Add(key string, db xorm.EngineInterface) {
	if _, ok := repo.database[key]; ok {
		return
	}
	repo.database[key] = db
}
