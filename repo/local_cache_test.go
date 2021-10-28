package repo

import (
	"encoding/gob"
	"github.com/subosito/gotenv"
	"testing"
	"time"
)

func TestGetLocalCacheRepo(t *testing.T) {
	if err := gotenv.Load("../.env"); err != nil {
		t.Error(err)
	}
	var repo = GetLocalCacheRepo()
	if repo == nil {
		t.Error("error cache open")
	}
}

func init() {
	gob.Register(&Options{})
	gob.Register(time.Time{})
}

func TestLocalCacheRepository_Set(t *testing.T) {
	if err := gotenv.Load("../.env"); err != nil {
		t.Error(err)
	}
	var (
		repo = GetLocalCacheRepo()
	)
	err := repo.SetMust(`test_options`,repo.options)
	if err != nil {
		t.Error(err)
	}

	defer repo.Sync()
}
