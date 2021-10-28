package entity

import (
	"github.com/weblfe/queue_mgr/utils"
)

type Properties KvMap

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
	var kv = KvMap(*p)
	return kv.Keys()
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

func (p *Properties) GetOr(key string, v ...interface{}) interface{} {
	var kv = KvMap(*p)
	return kv.Get(key, v...)
}

func (p *Properties) Exists(key string) bool {
	var kv = KvMap(*p)
	return kv.Exists(key)
}

func (p *Properties) VisitAll(each func(k string, v interface{})) {
	var (
		kv   = KvMap(*p)
		keys = p.Keys()
	)
	for _, k := range keys {
		each(k, kv[k])
	}
}

func (p *Properties) VisitAllCond(each func(k string, v interface{}) bool) {
	var (
		kv   = KvMap(*p)
		keys = p.Keys()
	)
	for _, k := range keys {
		if !each(k, kv[k]) {
			break
		}
	}
}
