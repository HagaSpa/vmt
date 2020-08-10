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

/*
Writer for Binary Operator (ARITHMETIC)

e.g.. add

1. pop to M register, decrease stack pointer by one.
	- @SP
	- M=M-1

2. push to D register from M register.
	- A=M
	- D=M

3. pop to M register, decrease stack pointer by one.
	- @SP
	- M=M-1

4. Add M register and D register, push to M register.
	- A=M
	- M=D+M

5. increase stack pointer by one.（Initialize stack pointer）
	- @SP
	- M=M+1

*/
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

// Writer for Unary Operator (ARITHMETIC)
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

// Writer for Condition Operator (ARITHMETIC)
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
