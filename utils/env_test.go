package utils

import (
	"os"
	"strings"
	"testing"
)

func TestParseEnvValue(t *testing.T) {
	var (
		v       = `default`
		err     = os.Setenv(`APP_NAME`, v)
		key     = `${APP_NAME}/dir/${ENV}`
		value   = ParseEnvValue(key)
		realStr = strings.Replace(value, `${APP_NAME}`, v, -1)
	)
	if err != nil {
		t.Error(err)
	}
	if value != realStr {
		t.Error("parse env failed ")
	}
}
