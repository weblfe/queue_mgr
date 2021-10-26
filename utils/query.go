package utils

import (
	"strings"
	"xorm.io/builder"
)

// CreateStrMultiCond 创建多条件查询条件builder
func CreateStrMultiCond(cons, key string) builder.Cond {
	if strings.Contains(cons, ",") {
		var arr = Str2Arr(cons)
		if len(arr) <= 0 {
			return nil
		}
		return builder.In(key, Str2Arr(cons))
	}
	return builder.Eq{key: cons}
}

// CreateStrMultiCondWithCheck 创建带检查的多条件查询条件builder
func CreateStrMultiCondWithCheck(cons, key string, check func(v string) bool) builder.Cond {
	if strings.Contains(cons, ",") {
		var arr = Str2Arr(cons)
		if len(arr) <= 0 {
			return nil
		}
		var arrIn []string
		for _, v := range arr {
			if check(v) {
				arrIn = append(arrIn, v)
			}
		}
		if len(arrIn) <= 0 {
			return nil
		}
		return builder.In(key, arrIn)
	}
	if !check(cons) {
		return nil
	}
	return builder.Eq{key: cons}
}

func CreateAndBuilder(cond builder.Cond, and ...builder.Cond) builder.Cond {
	if len(and) > 0 {
		for _, v := range and {
			if v == nil {
				continue
			}
			cond = v.And(cond)
		}
	}
	return cond
}

func CreateOrBuilder(cond builder.Cond, or ...builder.Cond) builder.Cond {
	if len(or) > 0 {
		for _, v := range or {
			cond = v.Or(cond)
		}
	}
	return cond
}
