package utils

import (
	rands "crypto/rand"
	"math"
	"math/big"
	"time"
)

// Random 随机
func Random(min, max int64) int64 {
	if min > max {
		min, max = max, min
	}
	if min < 0 {
		f64Min := math.Abs(float64(min))
		i64Min := int64(f64Min)
		result, _ := rands.Int(rands.Reader, big.NewInt(max+1+i64Min))
		return result.Int64() - i64Min
	}
	result, _ := rands.Int(rands.Reader, big.NewInt(max-min+1))
	return min + result.Int64()
}

// RandomInt 随机
func RandomInt(min, max int) int {
	return int(Random(int64(min), int64(max)))
}

// RandomSec 随机秒数
func RandomSec(min, max int64) time.Duration {
	return time.Duration(Random(min, max)) * time.Second
}
