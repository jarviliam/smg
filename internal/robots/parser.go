package robots

import "strings"

type Parser struct {
	tokens   []string
	position int
}

func NewParser(tokens []string) *Parser {
	p := &Parser{
		tokens: tokens,
	}
	return p
}

func (p *Parser) parseTokens() {

}

func (p *Parser) parseLine() {
	token1, ok1 := p.popToken()
	if !ok1 {
		return
	}
	//token2, ok2 := p.peekToken()
	// if !ok2 {
	// 	return
	// }

	switch strings.ToLower(token1) {

	}
}

func (p *Parser) popToken() (token string, ok bool) {
	token, ok = p.peekToken()
	if !ok {
		return
	}
	p.position++
	return token, true
}

func (p *Parser) peekToken() (token string, ok bool) {
	if p.position >= len(p.tokens) {
		return "", false
	}
	return p.tokens[p.position], true
}
