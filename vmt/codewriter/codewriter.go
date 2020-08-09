package codewriter

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
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

// Write for Unary Operator (ARITHMETIC)
func (cw *CodeWriter) writeUnaryOperator(op string) {
	asm :=
		`
// Unary Operator %s
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

func (cw *CodeWriter) writeConditionOperator(op string, n int) {
	var jump string
	switch op {
	case "eq":
		jump = "JEQ"
	case "gt":
		jump = "JGT"
	case "lt":
		jump = "JLT"
	}
	s := strconv.Itoa(n)
	asm :=
		`
// Condition Operator %s
@SP
M=M-1
A=M
D=M
@SP
M=M-1
A=M
D=M-D
@SP
A=M
M=-1
@%s
D;%s
@SP
A=M
M=0
(%s)
@SP
M=M+1
`
	asm = fmt.Sprintf(asm, op, s, jump, s)
	w := bufio.NewWriter(cw.w)
	w.WriteString(asm)
	w.Flush()
}
