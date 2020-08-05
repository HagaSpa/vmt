package parser

import (
	"bufio"
	"io"
	"strings"
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
	if line == "" || strings.HasPrefix(line, "//") {
		return
	}
	if strings.Contains(line, "//") {
		line = line[:strings.Index(line, "//")]
	}
	line = strings.TrimSpace(line)
	p.input = line
}
