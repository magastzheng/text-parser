package parser

import "github.com/mikydna/text-parser/lexer/lexeme"
import "github.com/mikydna/text-parser/parser/token"

type Parser struct {
	//input chan lexeme.Lexer
	resolvers []Resolver
	handlers []handler // could use map instead, but its a bit wacky
	bus chan []token.Token
}

type Resolver struct {
	accept lexeme.Predicate
	parse func(lexeme.Lexeme) []token.Token // TODO: needs err handling
}

type handler struct {
	accept token.Predicate
	process func(token.Token)
}

func New(resolvers ...Resolver) *Parser {
	return &Parser {
		//input: nil,
		resolvers: resolvers,
		handlers: []handler(nil)[:0],
		bus: make(chan []token.Token),
	}
}

func (par *Parser) Parse(input chan lexeme.Lexeme) {
	
	// how do i check if chan is closed, without pulling val? red herring?
	if input == nil { 
		return
	}

	// handle resolved lexeme objects
	go func() {
		for tokens := range par.bus {
			// tstructs in array could be mixed?
			for _, token := range tokens { 
				for _, handler := range par.handlers {
					if handler.accept(token) {
						go handler.process(token)
					}
				}
			}
		}
	}()
	

	// resolve incoming lexemes
	for lexeme := range input {
		for _ , resolver := range par.resolvers {
			if resolver.accept(lexeme) {
				// TODO: handle resolver error
				par.bus <- resolver.parse(lexeme)
			}
		} 
	}

	// need exit block; re: par.bus is closed ?
}