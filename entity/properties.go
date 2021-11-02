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
	var keys []string
	for k := range *p {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func (p *Properties) String() string {
	return utils.JsonEncode(p).String()
}

func (p *Properties) Len() int {
	if p == nil {
		return 0
	}
	return len(*p)
}

func (p *Properties) MarshalJSON() ([]byte, error) {
	var encoder = utils.JsonEncode(p)
	return encoder.Bytes(), encoder.Error()
}

func (p *Properties) UnmarshalJSON(data []byte) error {
	return utils.JsonDecode(data, p)
}

func (p *Properties) GetOr(key string, v ...string) string {
	if value, ok := (*p)[key]; ok {
		return utils.ParseEnvValue(value)
	}
	v = append(v, "")
	if v[0] != "" {
		return utils.ParseEnvValue(v[0])
	}
	return v[0]
}

func (p *Properties) Exists(key string) bool {
	if p == nil {
		return false
	}
	if _, ok := (*p)[key]; ok {
		return ok
	}
	return false
}

func (p *Properties) VisitAll(each func(k string, v interface{})) {
	var keys = p.Keys()
	for _, k := range keys {
		each(k, (*p)[k])
	}
}

func (p *Properties) Filters(each func(k string, v interface{}) bool) {
	var keys = p.Keys()
	for _, k := range keys {
		if !each(k, (*p)[k]) {
			break
		}
	}
}
