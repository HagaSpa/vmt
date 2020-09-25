package codewriter

import (
	"bufio"
	"fmt"
	"io"
	"vmt/parser"
)

type CodeWriter struct {
	w       io.Writer
	addr    int
	callcnt int
	fn      string
}

func New(w io.Writer) *CodeWriter {
	cw := &CodeWriter{
		w: w,
	}
	cw.writeInit()
	return cw
}

func (cw *CodeWriter) SetFileName(fn string) {
	cw.fn = fn
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
		cw.writeUnaryOperator("M=-M")
	case "not":
		cw.writeUnaryOperator("M=!M")
	case "eq":
		cw.writeConditionOperator("JEQ")
	case "gt":
		cw.writeConditionOperator("JGT")
	case "lt":
		cw.writeConditionOperator("JLT")
	}
}

func (cw *CodeWriter) WritePushPop(cmd parser.Type, segment string, index int) {
	switch cmd {
	case parser.PUSH:
		switch segment {
		case "constant":
			cw.writePushConstant(index)
		case "local":
			cw.writePushSymbol("LCL", index)
		case "argument":
			cw.writePushSymbol("ARG", index)
		case "this":
			cw.writePushSymbol("THIS", index)
		case "that":
			cw.writePushSymbol("THAT", index)
		case "pointer":
			cw.writePushRegister(index + 3)
		case "temp":
			cw.writePushRegister(index + 5)
		case "static":
			cw.writePushStatic(index)
		}
	case parser.POP:
		switch segment {
		case "local":
			cw.writePopSymbol("LCL", index)
		case "argument":
			cw.writePopSymbol("ARG", index)
		case "this":
			cw.writePopSymbol("THIS", index)
		case "that":
			cw.writePopSymbol("THAT", index)
		case "pointer":
			cw.writePopRegister(index + 3)
		case "temp":
			cw.writePopRegister(index + 5)
		case "static":
			cw.writePopStatic(index)
		}
	}
}

func (cw *CodeWriter) WriteLabel(label string) {
	symbol := fmt.Sprintf("%s$%s", cw.fn, label)
	asm := `
// write label %s
(%s)
`
	asm = fmt.Sprintf(asm, symbol, symbol)
	w := bufio.NewWriter(cw.w)
	w.WriteString(asm)
	w.Flush()
}

func (cw *CodeWriter) WriteIf(label string) {
	symbol := fmt.Sprintf("%s$%s", cw.fn, label)
	asm := `
// if-goto label %s
@SP
M=M-1
@SP
A=M
D=M
@%s
D;JNE
`
	asm = fmt.Sprintf(asm, symbol, symbol)
	w := bufio.NewWriter(cw.w)
	w.WriteString(asm)
	w.Flush()
}

func (cw *CodeWriter) WriteGoto(label string) {
	symbol := fmt.Sprintf("%s$%s", cw.fn, label)
	asm := `
// goto label %s
@%s
0;JMP
`
	asm = fmt.Sprintf(asm, symbol, symbol)
	w := bufio.NewWriter(cw.w)
	w.WriteString(asm)
	w.Flush()
}

func (cw *CodeWriter) WriteReturn() {
	asm := `
// return
@LCL
D=M
@R13
M=D
@5
D=A
@R13
A=M-D
D=M
@R14
M=D
@SP
M=M-1
@SP
A=M
D=M
@ARG
A=M
M=D
@ARG
D=M+1
@SP
M=D
@1
D=A
@R13
A=M-D
D=M
@THAT
M=D
@2
D=A
@R13
A=M-D
D=M
@THIS
M=D
@3
D=A
@R13
A=M-D
D=M
@ARG
M=D
@4
D=A
@R13
A=M-D
D=M
@LCL
M=D
@R14
A=M
0;JMP
`
	w := bufio.NewWriter(cw.w)
	w.WriteString(asm)
	w.Flush()
}

func (cw *CodeWriter) WriteFunction(funcname string, numlocal int) {
	asm := `
// function %s local nums %d
(%s)
`
	asm = fmt.Sprintf(asm, funcname, numlocal, funcname)
	w := bufio.NewWriter(cw.w)
	w.WriteString(asm)
	w.Flush()
	// initialize local variable to 0
	for i := 0; i < numlocal; i++ {
		cw.writePushConstant(0)
	}
}

func (cw *CodeWriter) WriteCall(funcname string, numargs int) {
	rlabel := fmt.Sprintf("%s$%d", funcname, cw.callcnt)
	asm := `
// call %s args nums %d
@%s
D=A
@SP
A=M
M=D
@SP
M=M+1
`
	asm = fmt.Sprintf(asm, funcname, numargs, rlabel)
	w := bufio.NewWriter(cw.w)
	w.WriteString(asm)
	w.Flush()

	// push LCL, ARG, THIS, THAT
	cw.writePushRegisterByName("LCL")
	cw.writePushRegisterByName("ARG")
	cw.writePushRegisterByName("THIS")
	cw.writePushRegisterByName("THAT")

	rasm := `
// set to return address 
@SP
D=M
@%d
D=D-A
@5
D=D-A
@ARG
M=D
@SP
D=M
@LCL
M=D
@%s
0;JMP
(%s)
`
	rasm = fmt.Sprintf(rasm, numargs, funcname, rlabel)
	w = bufio.NewWriter(cw.w)
	w.WriteString(rasm)
	w.Flush()

	cw.callcnt++
}

func (cw *CodeWriter) writeInit() {
	asm := `
// initialize asm
@256
D=A
@SP
M=D
`
	w := bufio.NewWriter(cw.w)
	w.WriteString(asm)
	w.Flush()
	cw.WriteCall("Sys.init", 0)
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
	asm := `
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
	- M=-M

3. increase stack pointer by one.（Initialize stack pointer）
	- @SP
	- M=M+1

*/
func (cw *CodeWriter) writeUnaryOperator(op string) {
	asm := `
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

6. if D register is 0, jump to LABEL1
	- @LABEL1
	- D;JEQ

7. Set FALSE to the memory pointed to by the stack pointer
	- @SP
	- A=M
	- M=0 //0x0000

8. Jump destination label
	- (LABEL1)

9. increase stack pointer by one.（Initialize stack pointer）
	- @SP
	- M=M+1

*/
func (cw *CodeWriter) writeConditionOperator(op string) {
	cw.addr++
	asm := `
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
@LABEL%d
D;%s
@SP
A=M
M=0
(LABEL%d)
@SP
M=M+1
`
	asm = fmt.Sprintf(asm, op, cw.addr, op, cw.addr)
	w := bufio.NewWriter(cw.w)
	w.WriteString(asm)
	w.Flush()
}

/*
Writer for Push Constant (PUSH)

e.g.. push constant 0

1. put 0 in D register
	- @0
	- D=A

2. pop to M register from D register. M register is top +1 element in stack
	- @SP
	- M=A //Set empty value, using stack pointer
	- M=D

3. increase stack pointer by one.（Initialize stack pointer）
	- @SP
	- M=M+1

*/
func (cw *CodeWriter) writePushConstant(index int) {
	asm := `
// push constant %d
@%d
D=A
@SP
A=M
M=D
@SP
M=M+1
`
	asm = fmt.Sprintf(asm, index, index)
	w := bufio.NewWriter(cw.w)
	w.WriteString(asm)
	w.Flush()
}

/*
Writer for Push Symbol (PUSH)

e.g.. push local 0

1. put 0 in D register
	- @0
	- D=A

2. Add the stack area pointed to by LCL(local) and the D register
	- @LCL
	- D=D+M

3. Save D register in the stack area pointed to by R13(temp)
	- @R13
	- M=D

4. Set the stack area pointed to by R13 in the A register
	- A=M
	- D=M

5. Save D register in the stack area pointed to by stack pointer
	- @SP
	- A=M
	- M=D

6. increase stack pointer by one.（Initialize stack pointer）
	- @SP
	- M=M+1

*/
func (cw *CodeWriter) writePushSymbol(symbol string, index int) {
	asm := `
// push symbol %s index %d 
@%d
D=A
@%s
D=D+M
@R13
M=D
A=M
D=M
@SP
A=M
M=D
@SP
M=M+1
`
	asm = fmt.Sprintf(asm, symbol, index, index, symbol)
	w := bufio.NewWriter(cw.w)
	w.WriteString(asm)
	w.Flush()
}

/*
Writer for Push Register (PUSH)

Set Register value in stack area ponted to by stack pointer.

e.g.. push temp 3

1. Put R8 Register value in D register. R8 is 5 + index(=3).
	- @R8
	- D=M

2. Add D register in the stack area pointed to by stack pointer
	- @SP
	- A=M
	- M=D

3. increase stack pointer by one.（Initialize stack pointer）
	- @SP
	- M=M+1

*/
func (cw *CodeWriter) writePushRegister(number int) {
	asm := `
// push register R%d
@R%d
D=M
@SP
A=M
M=D
@SP
M=M+1
`
	asm = fmt.Sprintf(asm, number, number)
	w := bufio.NewWriter(cw.w)
	w.WriteString(asm)
	w.Flush()
}

func (cw *CodeWriter) writePushRegisterByName(register string) {
	asm := `
// push register %s
@%s
D=M
@SP
A=M
M=D
@SP
M=M+1
`
	asm = fmt.Sprintf(asm, register, register)
	w := bufio.NewWriter(cw.w)
	w.WriteString(asm)
	w.Flush()
}

/*
Writer for Push Static (PUSH)

Set Static value in stack area ponted to by stack pointer.

e.g.. push static 1 (StackTest.vm)

1. Put Static value in D register.
	- @StaticTest.1
	- D=M

2. Add D register in the stack area pointed to by stack pointer
	- @SP
	- A=M
	- M=D

3. increase stack pointer by one.（Initialize stack pointer）
	- @SP
	- M=M+1

*/
func (cw *CodeWriter) writePushStatic(index int) {
	static := fmt.Sprintf("%s.%d", cw.fn, index)
	asm := `
// push static %s
@%s
D=M
@SP
A=M
M=D
@SP
M=M+1
`
	asm = fmt.Sprintf(asm, static, static)
	w := bufio.NewWriter(cw.w)
	w.WriteString(asm)
	w.Flush()
}

/*
Writer for Pop Symbol (POP)

Pop the data at the top of the stack and store it in segment[index]

e.g.. pop local 0

1. pop to M register, decrease stack pointer by one.
	- @SP
	- M=M-1

2. put 0 in D register
	- @0
	- D=A

3. Add the stack area pointed to by LCL(local) and the D register
	- @LCL
	- D=D+M

4. put the data at the top of the stack in D register
	- @SP
	- A=M
	- D=M

5. Add D register in the stack area pointed to by R13(temp)
	- @R13
	- A=M
	- M=D

*/
func (cw *CodeWriter) writePopSymbol(symbol string, index int) {
	asm := `
// pop symbol %s index %d
@SP
M=M-1
@%d
D=A
@%s
D=D+M
@R13
M=D
@SP
A=M
D=M
@R13
A=M
M=D
`
	asm = fmt.Sprintf(asm, symbol, index, index, symbol)
	w := bufio.NewWriter(cw.w)
	w.WriteString(asm)
	w.Flush()
}

/*
Writer for Pop Register (POP)

Pop the data at the top of the stack and store it in register

e.g.. pop temp 6

1. pop to M register, decrease stack pointer by one.
	- @SP
	- M=M-1

4. put the data at the top of the stack in D register
	- @SP
	- A=M
	- D=M

5. Add D register in the stack area pointed to by R11. R11 is 5 + index(=6).
	- @R11
	- M=D

*/
func (cw *CodeWriter) writePopRegister(index int) {
	asm := `
// pop register R%d
@SP
M=M-1
@SP
A=M
D=M
@R%d
M=D
`
	asm = fmt.Sprintf(asm, index, index)
	w := bufio.NewWriter(cw.w)
	w.WriteString(asm)
	w.Flush()
}

/*
Writer for Pop Static (POP)

Pop the data at the top of the stack and store it in static area (Memory[static])

e.g.. pop static 0 (StackTest.vm)

1. pop to M register, decrease stack pointer by one.
	- @SP
	- M=M-1

2. put the data at the top of the stack in D register
	- @SP
	- A=M
	- D=M

3. Add D register in the stack area pointed to by static.
	- @StaticTest.0
	- M=D

*/
func (cw *CodeWriter) writePopStatic(index int) {
	static := fmt.Sprintf("%s.%d", cw.fn, index)
	asm := `
// pop static %s
@SP
M=M-1
@SP
A=M
D=M
@%s
M=D
`
	asm = fmt.Sprintf(asm, static, static)
	w := bufio.NewWriter(cw.w)
	w.WriteString(asm)
	w.Flush()
}
