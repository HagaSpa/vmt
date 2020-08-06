package parser

import (
	"bufio"
	"io"
	"strings"
)

type Type int

// command type
const (
	_ Type = iota
	None
	ARITHMETIC
	PUSH
	POP
	LABEL
	GOTO
	IF
	FUNCTION
	RETURN
	CALL
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

func (p *Parser) CommandType() Type {
	slice := strings.Split(p.input, " ")
	switch slice[0] {
	case "add", "sub", "neg", "eq", "gt", "lt", "and", "or", "not":
		return ARITHMETIC
	case "push":
		return PUSH
	case "pop":
		return POP
	case "label":
		return LABEL
	case "goto":
		return GOTO
	case "if-goto":
		return IF
	case "function":
		return FUNCTION
	case "return":
		return RETURN
	case "call":
		return CALL
	default:
		return None // invalid vm command
	}
}

func (p *Parser) Arg1() string {
	slice := strings.Split(p.input, " ")
	switch slice[0] {
	case "add", "sub", "neg", "eq", "gt", "lt", "and", "or", "not":
		return slice[0]
	default:
		return slice[1]
	}
}
