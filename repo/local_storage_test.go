package repo

import (
	"fmt"
	"github.com/subosito/gotenv"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"testing"
)

func TestLocalStorageRepository_GetStorage(t *testing.T) {
	if err := gotenv.Load("../.env"); err != nil {
		t.Error(err)
	}
	var repo = GetLocalStorageRepo()
	if db, ok := repo.GetStorage(); !ok || db == nil {
		t.Error("error db open")
	}
}

func TestLocalStorageRepository_Put(t *testing.T) {
	if err := gotenv.Load("../.env"); err != nil {
		t.Error(err)
	}
	var (
		repo   = GetLocalStorageRepo()
		db, ok = repo.GetStorage()
	)
	if !ok || db == nil {
		t.Error("error db open")
	}
	v, err := db.Get([]byte(`test`), &opt.ReadOptions{})
	if err != nil {
		t.Error(err)
	}
	fmt.Println(string(v))
}
