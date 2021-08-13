package robots

import (
	"bytes"
	"go/token"
	"sync"
	"unicode/utf8"
)

type Reader struct {
	data     []byte
	position token.Position
	current  rune
}

const EOL = "\n"

var tokBuffers = sync.Pool{New: func() interface{} { return bytes.NewBuffer(make([]byte, 32)) }}

func NewReader(data []byte) *Reader {
	r := &Reader{
		data: data,
	}
	return r
}

func (r *Reader) ExtractTokens() []string {
	r.position.Offset = 0
	r.position.Line = 1
	r.position.Column = 1
	output := make([]string, 0, 64)
	for {
		token := r.readData()
		if token == "" {
			break
		}
		output = append(output, token)
	}
	return output
}

func (r *Reader) readData() string {
	if r.isOverflow() {
		return ""
	}

	r.skipIfSpace()
	if r.current == -1 {
		return ""
	}
	//TODO: Handle Comments
	if r.current == '#' {
		r.skipTillEOL()
		if r.current == -1 {
			return ""
		}
		return EOL
	}
	//Parse Token (Until Space / Newline)
	token := tokBuffers.Get().(*bytes.Buffer)
	defer tokBuffers.Put(token)
	token.Reset()
	token.WriteRune(r.current)
	r.readCharacter()
	for r.current != -1 && !r.isCurrentSpace() && !r.isCurrentNewLine() {
		token.WriteRune(r.current)
		r.readCharacter()
	}
	return token.String()
}

/*
 * Read Characters
 */
func (r *Reader) readCharacter() bool {
	if r.isOverflow() {
		r.current = -1
		return false
	}
	r.position.Column += 1
	if r.isCurrentNewLine() {
		//Move to next line
		r.position.Line++
		r.position.Column = 1
	}
	ru, w := rune(r.data[r.position.Offset]), 1
	if ru >= 0x80 {
		ru, w = utf8.DecodeRune(r.data[r.position.Offset:])
	}
	r.position.Column += 1
	r.position.Offset += w
	r.current = ru
	return true
}

//Utility Functions
func (r *Reader) isCurrentNewLine() bool {
	return r.current == '\n'
}
func (r *Reader) isCurrentSpace() bool {
	return r.current == ' '
}

func (r *Reader) isOverflow() bool {
	return r.position.Offset > len(r.data)
}

func (r *Reader) skipTillEOL() {
	for r.current != -1 && !r.isCurrentNewLine() {
		r.readCharacter()
	}

	for r.current != -1 && r.isCurrentNewLine() {
		r.readCharacter()
	}
}

func (r *Reader) skipIfSpace() {
	for r.current != -1 && r.current == ' ' {
		r.readCharacter()
	}
}
