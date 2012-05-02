package lexeme

import (
	"fmt"	
)

type Kind uint8

type Lexeme struct {
	Kind Kind
	Pos Position
	Part []byte
}

func (l *Lexeme) String() string {
	return fmt.Sprintf("Lexeme(%v, %v)", l.Kind, string(l.Part))
}