package repo

import (
	"github.com/subosito/gotenv"
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
