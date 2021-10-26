package repo

import "testing"

func TestRedisOf(t *testing.T) {
	var rd = RedisDb("push")

	rd.Lock("test")
}
