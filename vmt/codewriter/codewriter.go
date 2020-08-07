package codewriter

import (
	"bufio"
	"fmt"
	"io"
	"vmt/parser"
)

type CodeWriter struct {
	w io.Writer
	p *parser.Parser
}

func New(w io.Writer, p *parser.Parser) *CodeWriter {
	cw := &CodeWriter{
		w: w,
		p: p,
	}
	return cw
}

// Writer for Binary Operator (ARITHMETIC)
func (cw *CodeWriter) writeBinaryOperator(op string) {
	asm :=
		`
// Binary Operator %s
@SP
M=M-1
A=M
D=M
@SP
M=M-1
A=M
%s
@SP
M=M+1
`
	asm = fmt.Sprintf(asm, op, op)
	w := bufio.NewWriter(cw.w)
	w.WriteString(asm)
	w.Flush()
}
