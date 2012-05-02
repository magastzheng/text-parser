package lexeme

import (
	b "bytes"
)

import "github.com/mikydna/text-parser/bytes"

type Predicate func(Lexeme) bool

func OfKind(k Kind) Predicate {
	return func(t Lexeme) bool {
		return t.Kind == k
	}
}

func EmptyValue() Predicate {
	return func(l Lexeme) bool {
		trimmed := b.TrimSpace(l.Part)
		return len(trimmed) == 0 // is a zero-slice equal to nil?
	}
}

func ValueOf(p bytes.Predicate) Predicate {
	return func(l Lexeme) bool {
		ok, _ := p(l.Part)
		return ok
	}
}

func And(ps ...Predicate) Predicate {
	return func(l Lexeme) bool {
		result := len(ps) != 0
		for _, p := range ps {
			result = result && p(l)
			if !result {
				return false
			}
		}
		return result
	}
}

func Or(ps ...Predicate) Predicate {
	return func(l Lexeme) bool {
		result := false
		for _, p := range ps {
			result = result || p(l)
		}
		return result
	}
}
