package docs

import "fmt"

const (
	schemaDefault = "http:"
)

func GetDefaultSchema() string{
	return fmt.Sprintf("%s//",schemaDefault)
}