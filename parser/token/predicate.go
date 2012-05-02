package token

import (
	"reflect"
)

type Predicate func(Token) bool

func OfType(s reflect.Type) func(Token) bool {
  return func(t Token) bool {
    return s == reflect.TypeOf(t)
  }
}
