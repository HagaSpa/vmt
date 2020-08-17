package codewriter

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
)

type CodeWriter struct {
	w    io.Writer
	addr int
}

func New(w io.Writer) *CodeWriter {
	cw := &CodeWriter{
		w: w,
	}
	return cw
}

func (cw *CodeWriter) WriteArithmetic(cmd string) {
	switch cmd {
	case "add":
		cw.writeBinaryOperator("M=D+M")
	case "sub":
		cw.writeBinaryOperator("M=M-D")
	case "and":
		cw.writeBinaryOperator("M=D&M")
	case "or":
		cw.writeBinaryOperator("M=D|M")
	case "neg":
		cw.writeUnaryOperator("-M")
	case "not":
		cw.writeUnaryOperator("!M")
	case "eq":
		cw.writeConditionOperator("eq")
	case "gt":
		cw.writeConditionOperator("gt")
	case "lt":
		cw.writeConditionOperator("lt")
	}
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

/*
Writer for Unary Operator (ARITHMETIC)

e.g.. neg

1. pop to M register, decrease stack pointer by one.
	- @SP
	- M=M-1

2. Make it a negative number at M register.
	- A=M
	- -M

3. increase stack pointer by one.（Initialize stack pointer）
	- @SP
	- M=M+1

*/
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

/*
Writer for Condition Operator (ARITHMETIC)

e.g.. eq

1. pop to M register, decrease stack pointer by one.
	- @SP
	- M=M-1

2. push to D register from M register.
	- A=M
	- D=M

3. pop to M register, decrease stack pointer by one.
	- @SP
	- M=M-1

4. substract M register from D register, push to D register.
	- A=M
	- D=M-D

5. Set TRUE to the memory pointed to by the stack pointer
	- @SP
	- A=M
	- M=-1 //0xFFFF

6. if D register is 0, jump to 1
	- @1
	- D;JEQ

7. Set FALSE to the memory pointed to by the stack pointer
	- @SP
	- A=M
	- M=0 //0x0000

8. Jump destination label
	- (1)

9. increase stack pointer by one.（Initialize stack pointer）
	- @SP
	- M=M+1

*/
func (cw *CodeWriter) writeConditionOperator(op string) {
	var jump string
	switch op {
	case "eq":
		jump = "JEQ"
	case "gt":
		jump = "JGT"
	case "lt":
		jump = "JLT"
	}
	cw.addr++
	s := strconv.Itoa(cw.addr)
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

func (cw *CodeWriter) writePushConstant(index int) {
	s := strconv.Itoa(index)
	asm := `
// push constant %s
@%s
D=A
@SP
A=M
M=D
@SP
M=M+1
`
	asm = fmt.Sprintf(asm, s, s)
	w := bufio.NewWriter(cw.w)
	w.WriteString(asm)
	w.Flush()
}
