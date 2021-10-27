package entity

type (
	FastCgiType string
)

func (t FastCgiType) String() string {
	return string(t)
}

func (t FastCgiType) Eq(v string) bool {
	return string(t) == v
}

func (t FastCgiType) Equal(v FastCgiType) bool {
	return t == v
}
