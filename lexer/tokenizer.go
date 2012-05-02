package lexer

import "github.com/mikydna/text-parser/bytes"
import "github.com/mikydna/text-parser/lexer/lexeme"

type Tokenizer interface {
	match([]byte) bool
	consume(*Lexer) Tokenizer
}

type DelegatingT struct {
	Kind lexeme.Kind
}

// --- impl

func (_ *DelegatingT) match(b []byte) bool {
	return b != nil
}

func (t *DelegatingT) consume(lex *Lexer) Tokenizer {
	for {
		for _, lexT := range lex.lexTs {
			if lexT.match(lex.input[lex.cursor:]) {
				pos, b := lex.capture()
				lex.publish(lexeme.Lexeme{t.Kind, pos, b})
			}
			if lex.next() == EOF {
				pos, b := lex.capture()
				lex.publish(lexeme.Lexeme{t.Kind, pos, b})
				break
			}
		}
	}
	return nil
}

// --- impl

type PredicateT struct {
	Kind lexeme.Kind
	PfxPred bytes.Predicate
	SfxPred bytes.Predicate
}

func (t *PredicateT) match(b []byte) bool {
	ok, _ := t.PfxPred(b)
	return ok
}

func (t *PredicateT) consume(lex *Lexer) Tokenizer {
	for {
		ok, _ := t.SfxPred(lex.input[lex.anchor:lex.cursor])
		if ok {
			pos, b := lex.capture()
			lex.publish(lexeme.Lexeme{t.Kind, pos, b})
			return lex.baseT
		}
		if lex.next() == EOF {
			pos, b := lex.capture()
			lex.publish(lexeme.Lexeme{t.Kind, pos, b})
			break
		}
	}
	return nil
}