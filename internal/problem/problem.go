// Package problem carries application-owned message keys and values to the UI.
// Values stay separate so localization never rewrites user data.
package problem

import (
	"fmt"
	"regexp"
)

type Error struct {
	Key    string
	Values map[string]any
}

func New(key string, values map[string]any) *Error { return &Error{Key: key, Values: values} }

var placeholder = regexp.MustCompile(`\{(\w+)\}`)

func (e *Error) Error() string {
	return placeholder.ReplaceAllStringFunc(e.Key, func(token string) string {
		if value, ok := e.Values[token[1:len(token)-1]]; ok {
			return fmt.Sprint(value)
		}
		return token
	})
}
