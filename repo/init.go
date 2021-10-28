package repo

import (
	"encoding/gob"
	"time"
)

func init() {
	gob.Register(&Options{})
	gob.Register(&time.Time{})
}
