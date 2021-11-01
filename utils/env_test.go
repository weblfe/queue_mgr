package utils

import (
	"os"
	"strings"
	"testing"
)

func TestParseEnvValue(t *testing.T) {
	var (
		v       = `default`
		_       = os.Setenv(`APP_NAME`, v)
		key     = `${APP_NAME}/dir/${ENV}`
		value   = ParseEnvValue(key)
		realStr = strings.Replace(value, `${APP_NAME}`, v, -1)
	)

	if value != realStr {
		t.Error("parse env failed ")
	}
}
