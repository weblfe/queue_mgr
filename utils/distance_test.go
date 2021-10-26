package utils

import (
	"fmt"
	"testing"
)

func TestPoint_GetDistance(t *testing.T) {
	var (
		p1 = Point{
			1.00, 2.00,
		}
		p2 = Point{
			10.00, 60.00,
		}
		distance =  p1.GetDistance(p2)
	)
	if distance < 0 {
		t.Error("计算异常")
	}
	fmt.Println(distance)
}
