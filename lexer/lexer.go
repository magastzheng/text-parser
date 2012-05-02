package lexer

import (
	"unicode/utf8"
)

import "github.com/mikydna/text-parser/lexer/lexeme"

const (
	_ rune = -(iota + 1)
	EOF
)

type Lexer struct {
	input []byte
	anchor, cursor int
	baseT Tokenizer 
	lexTs []Tokenizer
	handlers []handler
	bus chan lexeme.Lexeme
}

type handler struct {
	accept lexeme.Predicate
	process func(lexeme.Lexeme) //needs error handling
}

func New(base Tokenizer, lexs ...Tokenizer) *Lexer {
	return &Lexer {
		input: []byte(nil)[:0],
		anchor: 0,
		cursor: 0,
		baseT: base,
		lexTs: lexs,
		handlers: []handler(nil)[:0],
		bus: make(chan lexeme.Lexeme),
	}
}

func (lex *Lexer) Scan(input []byte) {
	lex.input = input
	lex.anchor = 0
	lex.cursor = 0

	// token handling
	go func() {
		for lexeme := range lex.bus {
			for _, handler := range lex.handlers {
				if handler.accept(lexeme) {
					go handler.process(lexeme)
				}
			}
		}
	}()

	// process text
	var t Tokenizer = lex.baseT
	for t != nil {
		t = t.consume(lex)
	}
}

func (lex *Lexer) Handle(pred lexeme.Predicate, action func(lexeme.Lexeme)) {
	lex.handlers = append(lex.handlers, handler{pred,action})
}

func (lex *Lexer) next() rune {
	if lex.cursor >= len(lex.input) {
		return EOF
	}
	r, w := utf8.DecodeRune(lex.input[lex.cursor:])
	lex.cursor += w
	return r
}

func (lex *Lexer) capture() (lexeme.Position, []byte) {
	pos := lexeme.PositionOf(lex.input, lex.anchor)
	if v := lex.input[lex.anchor:lex.cursor]; len(v) > 0 {
		lex.anchor = lex.cursor
		return *pos, v
	}
	return *pos, nil
}

func (lex *Lexer) publish(l lexeme.Lexeme) {
	lex.bus <- l
}