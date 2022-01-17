package pkg

import "fmt"

type ContextKey string

func (c ContextKey) String() string {
	return fmt.Sprintf("%s%s", contextKeyPrefix, string(c))
}

const (
	contextKeyPrefix = "goDevopsLogging-"
)
