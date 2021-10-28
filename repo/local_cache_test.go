package repo

import (
	"github.com/subosito/gotenv"
	"testing"
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

func TestLocalCacheRepository_Set(t *testing.T) {
	if err := gotenv.Load("../.env"); err != nil {
		t.Error(err)
	}
	var (
		repo = GetLocalCacheRepo()
	)
	err := repo.SetMust(`test_options`, repo.options)
	if err != nil {
		t.Error(err)
	}
	err = repo.SetMust("hashMap", map[string]interface{}{
		"number":   1,
		"optional": repo.options,
	})
	if err != nil {
		t.Error(err)
	}
	defer repo.Sync()
}
