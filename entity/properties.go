package entity

import (
	"github.com/weblfe/queue_mgr/utils"
	"sort"
)

type Properties map[string]string

func NewProperties() *Properties {
	var p = new(Properties)
	return p
}

func ParseProperties(data string) (*Properties, error) {
	var p = NewProperties()
	if err := p.UnmarshalJSON([]byte(data)); err != nil {
		return nil, err
	}
	return p, nil
}

func (p *Properties) Keys() []string {
	var (
		keys []string
		kv   = map[string]string(*p)
	)
	for k := range kv {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func (p *Properties) String() string {
	return utils.JsonEncode(p).String()
}

func (p *Properties) MarshalJSON() ([]byte, error) {
	var encoder = utils.JsonEncode(p)
	return encoder.Bytes(), encoder.Error()
}

func (p *Properties) UnmarshalJSON(data []byte) error {
	return utils.JsonDecode(data, p)
}

func (p *Properties) GetOr(key string, v ...string) string {
	if v, ok := (*p)[key]; ok {
		return v
	}
	if len(v) <= 0 {
		v = append(v, "")
	}
	return v[0]
}

func (p *Properties) Exists(key string) bool {
	if _, ok := (*p)[key]; ok {
		return ok
	}
	return false
}

func (p *Properties) VisitAll(each func(k string, v interface{})) {
	var (
		kv   = map[string]string(*p)
		keys = p.Keys()
	)
	for _, k := range keys {
		each(k, kv[k])
	}
}

func (p *Properties) VisitAllCond(each func(k string, v interface{}) bool) {
	var (
		kv   = map[string]string(*p)
		keys = p.Keys()
	)
	for _, k := range keys {
		if !each(k, kv[k]) {
			break
		}
	}
}
