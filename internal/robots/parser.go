package robots

import (
	"go/token"
	"unicode/utf8"
)

type Parser struct {
	isLoaded bool
	position token.Position
	input    []byte
	curr     rune
}

func NewParser() *Parser {
	p := &Parser{}
	return p
}

/**
* Extracts The Robots Text File
* Returns a List of String Tokens
 */
func (p *Parser) Extract() []string {
	if !p.isLoaded {
		return
	}
	output := make([]string, 0, 64)
	for {
		token := p.read()
		if token == "" {
			break
		}
		output = append(output, token)
	}
	return output
}

/*
* Basic Parser Works as :
* 1) Check to see if not ovrflowing
* 2) Skip all Spaces
* 3) If The next non space rune is a new line, add to string array (new line is a token). Skip duplicate newlines
* 4) If there is a non comment rune then we have a token -> Tokenize?
* 5) Outcome should be an ordered list of tokens (Keywords / New lines)
 */
func (p *Parser) read() string {
	if p.isOverflow() {
		return ""
	}

	p.skipIfSpace()
	if p.curr == -1 {
		return ""
	}
	//TODO: Handle NewLine
	//TODO: Handle Comments
	if p.curr == '#' {

	}

	//Parse Token (Until Space / Newline)
}

func (p *Parser) skipIfSpace() {
	for p.curr != -1 && p.curr == ' ' {
	}
}
func (p *Parser) readChar() bool {
	if p.isOverflow() {
		p.curr = -1
		return false
	}
	p.position.Column += 1
	if p.isCurrentNewLine() {
		//Move to next line
		p.position.Line++
		p.position.Column = 1
	}
	ru, w := rune(p.input[p.position.Offset]), 1
	if ru >= 0x80 {
		ru, w = utf8.DecodeRune(p.input[p.position.Offset:])
	}
	p.position.Column += 1
	p.position.Offset += w
	p.curr = ru
	return true
}

func (p *Parser) isCurrentNewLine() bool {
	return p.curr == '\n'
}
func (p *Parser) isCurrentSpace() bool {
	return p.curr == ' '
}

func (p *Parser) isOverflow() bool {
	return p.position.Offset > len(p.input)
}
