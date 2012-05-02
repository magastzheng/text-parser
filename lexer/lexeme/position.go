package lexeme

import (
	"bytes"
)

var (
  NewLine []byte = []byte("\n")
)
        
type Position struct {
  Line, Column int
}

func PositionOf(doc []byte, offset int) *Position {
  splits := bytes.Split(doc[offset:], NewLine)
  nlines := len(splits)
  return &Position{nlines-1, len(splits[nlines-1])}
}