package parser

import (
	"bufio"
	"io"
)

type Parser struct {
	scanner *bufio.Scanner
	input   string
}

func New(r io.Reader) *Parser {
	s := bufio.NewScanner(r)
	p := &Parser{
		scanner: s,
	}
	return p
}

func (p *Parser) HasMoreCommands() bool {
	return p.scanner.Scan()
}

func (p *Parser) Advance() {
	line := p.scanner.Text()
	p.input = line
}
