package bytes

import (
	"bytes"
	"unicode/utf8"
)

// Bytes predicates evaluate the []byte as a sequence,
// and returns whether the []byte sequence can satisfy
// the predicate and at what array index did the match 
// attempt fail or finish at.

// A bytes/predicate returns a tuple (bool, int):
// 1. bool: does the byte[] satisfy the predicate, 
// 2. int: the match index

type Predicate func([]byte) (bool, int)

func Is(t rune) Predicate {
	return func(b []byte) (bool, int) {
		r, n := utf8.DecodeRune(b)
		if r == t {
			return true, n
		}
		return false, 0
	}
}

func Not(t rune) Predicate {
	return func(b []byte) (bool, int) {
		r, n := utf8.DecodeRune(b)
		if r != t {
			return true, n
		}
		return false, 0
	}
}

// ops

func Repeats(p Predicate) Predicate {
	return func(b []byte) (bool, int) {
		for i := range b {
			ok, _ := p(b[i:])
			if !ok {
				return i != 0, i
			}
		}
		return len(b) != 0, len(b)
	}
}

func Exact(p Predicate) Predicate {
	return func(b []byte) (bool, int) {
		ok, n := p(b)
		return ok && n == len(b), n
	}
}

// Reverses the input, then evals the predicate.
// TODO: Costly; Profile.
func Reverse(p Predicate) Predicate {
	return func(b []byte) (bool, int) {
		runes := bytes.Runes(b)
		length := len(runes)
		reversed := make([]rune, length)
		for i := range runes {
			reversed[length-1-i] = runes[i]
		}
		return p([]byte(string(reversed)))
	}
}

// builders

func (p Predicate) Then(q Predicate) Predicate {
	return func(b []byte) (bool, int) {
		pok, pn := p(b)
		if !pok {
			return false, pn
		}
		qok, qn := q(b[pn:])
		return qok, qn + pn
	}
}
