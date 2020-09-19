package codewriter

import (
	"bytes"
	"reflect"
	"testing"
	"vmt/parser"
)

func TestCodeWriter_writeBinaryOperator(t *testing.T) {
	tests := []struct {
		name string
		line string
		args string
		want string
	}{
		{
			"add",
			"add",
			"M=D+M",
			`
// Binary Operator M=D+M
@SP
M=M-1
A=M
D=M
@SP
M=M-1
A=M
M=D+M
@SP
M=M+1
`,
		},
		{
			"sub",
			"sub",
			"M=M-D",
			`
// Binary Operator M=M-D
@SP
M=M-1
A=M
D=M
@SP
M=M-1
A=M
M=M-D
@SP
M=M+1
`,
		},
		{
			"and",
			"and",
			"M=D&M",
			`
// Binary Operator M=D&M
@SP
M=M-1
A=M
D=M
@SP
M=M-1
A=M
M=D&M
@SP
M=M+1
`,
		},
		{
			"or",
			"or",
			"M=D|M",
			`
// Binary Operator M=D|M
@SP
M=M-1
A=M
D=M
@SP
M=M-1
A=M
M=D|M
@SP
M=M+1
`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := bytes.NewBufferString("")
			cw := &CodeWriter{
				w: b,
			}
			cw.writeBinaryOperator(tt.args)

			if string(b.Bytes()) != tt.want {
				t.Errorf("writeBinaryOperator() = %s, want %v", b, tt.want)
			}
		})
	}
}

func TestCodeWriter_writeUnaryOperator(t *testing.T) {
	tests := []struct {
		name string
		line string
		args string
		want string
	}{
		{
			"neg",
			"neg",
			"-M",
			`
// Unary Operator -M
@SP
M=M-1
A=M
-M
@SP
M=M+1
`,
		},
		{
			"not",
			"not",
			"!M",
			`
// Unary Operator !M
@SP
M=M-1
A=M
!M
@SP
M=M+1
`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := bytes.NewBufferString("")
			cw := &CodeWriter{
				w: b,
			}
			cw.writeUnaryOperator(tt.args)

			if string(b.Bytes()) != tt.want {
				t.Errorf("writeUnaryOperator() = %s, want %v", b, tt.want)
			}
		})
	}
}

func TestCodeWriter_writeConditionOperator(t *testing.T) {
	tests := []struct {
		name string
		op   string
		want string
	}{
		{
			"eq",
			"JEQ",
			`
// Condition Operator JEQ
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
@LABEL1
D;JEQ
@SP
A=M
M=0
(LABEL1)
@SP
M=M+1
`,
		},
		{
			"gt",
			"JGT",
			`
// Condition Operator JGT
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
@LABEL1
D;JGT
@SP
A=M
M=0
(LABEL1)
@SP
M=M+1
`,
		},
		{
			"lt",
			"JLT",
			`
// Condition Operator JLT
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
@LABEL1
D;JLT
@SP
A=M
M=0
(LABEL1)
@SP
M=M+1
`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := bytes.NewBufferString("")
			cw := &CodeWriter{
				w: b,
			}
			cw.writeConditionOperator(tt.op)

			if string(b.Bytes()) != tt.want {
				t.Errorf("writeConditionOperator() = %s, want %v", b, tt.want)
			}
		})
	}
}

func TestCodeWriter_writeConditionOperator_incrementAddr(t *testing.T) {
	tests := []struct {
		name string
		op   string
		want string
	}{
		{
			"eq",
			"JEQ",
			`
// Condition Operator JEQ
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
@LABEL1
D;JEQ
@SP
A=M
M=0
(LABEL1)
@SP
M=M+1
`,
		},
		{
			"gt",
			"JGT",
			`
// Condition Operator JEQ
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
@LABEL1
D;JEQ
@SP
A=M
M=0
(LABEL1)
@SP
M=M+1

// Condition Operator JGT
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
@LABEL2
D;JGT
@SP
A=M
M=0
(LABEL2)
@SP
M=M+1
`,
		},
		{
			"lt",
			"JLT",
			`
// Condition Operator JEQ
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
@LABEL1
D;JEQ
@SP
A=M
M=0
(LABEL1)
@SP
M=M+1

// Condition Operator JGT
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
@LABEL2
D;JGT
@SP
A=M
M=0
(LABEL2)
@SP
M=M+1

// Condition Operator JLT
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
@LABEL3
D;JLT
@SP
A=M
M=0
(LABEL3)
@SP
M=M+1
`,
		},
	}

	b := bytes.NewBufferString("")
	cw := &CodeWriter{
		w: b,
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cw.writeConditionOperator(tt.op)
			if string(b.Bytes()) != tt.want {
				t.Errorf("writeConditionOperator() = %s, want %v", b, tt.want)
			}
		})
	}
}

func TestCodeWriter_WriteArithmetic(t *testing.T) {
	tests := []struct {
		name string
		cmd  string
		want string
	}{
		{
			"add",
			"add",
			`
// Binary Operator M=D+M
@SP
M=M-1
A=M
D=M
@SP
M=M-1
A=M
M=D+M
@SP
M=M+1
`,
		},
		{
			"sub",
			"sub",
			`
// Binary Operator M=M-D
@SP
M=M-1
A=M
D=M
@SP
M=M-1
A=M
M=M-D
@SP
M=M+1
`,
		},
		{
			"and",
			"and",
			`
// Binary Operator M=D&M
@SP
M=M-1
A=M
D=M
@SP
M=M-1
A=M
M=D&M
@SP
M=M+1
`,
		},
		{
			"or",
			"or",
			`
// Binary Operator M=D|M
@SP
M=M-1
A=M
D=M
@SP
M=M-1
A=M
M=D|M
@SP
M=M+1
`,
		},
		{
			"neg",
			"neg",
			`
// Unary Operator M=-M
@SP
M=M-1
A=M
M=-M
@SP
M=M+1
`,
		},
		{
			"not",
			"not",
			`
// Unary Operator M=!M
@SP
M=M-1
A=M
M=!M
@SP
M=M+1
`,
		},
		{
			"eq",
			"eq",
			`
// Condition Operator JEQ
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
@LABEL1
D;JEQ
@SP
A=M
M=0
(LABEL1)
@SP
M=M+1
`,
		},
		{
			"gt",
			"gt",
			`
// Condition Operator JGT
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
@LABEL1
D;JGT
@SP
A=M
M=0
(LABEL1)
@SP
M=M+1
`,
		},
		{
			"lt",
			"lt",
			`
// Condition Operator JLT
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
@LABEL1
D;JLT
@SP
A=M
M=0
(LABEL1)
@SP
M=M+1
`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := bytes.NewBufferString("")
			cw := &CodeWriter{
				w: b,
			}
			cw.WriteArithmetic(tt.cmd)

			if string(b.Bytes()) != tt.want {
				t.Errorf("WriteArithmetic() = %s, want %v", b, tt.want)
			}
		})
	}
}

func TestCodeWriter_WritePushPop(t *testing.T) {
	type args struct {
		cmd     parser.Type
		segment string
		index   int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"push constant 100",
			args{
				cmd:     parser.PUSH,
				segment: "constant",
				index:   100,
			},
			`
// push constant 100
@100
D=A
@SP
A=M
M=D
@SP
M=M+1
`,
		},
		{
			"push local 0",
			args{
				cmd:     parser.PUSH,
				segment: "local",
				index:   0,
			},
			`
// push symbol LCL index 0 
@0
D=A
@LCL
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
`,
		},
		{
			"push argument 2",
			args{
				cmd:     parser.PUSH,
				segment: "argument",
				index:   2,
			},
			`
// push symbol ARG index 2 
@2
D=A
@ARG
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
`,
		},
		{
			"push this 6",
			args{
				cmd:     parser.PUSH,
				segment: "this",
				index:   6,
			},
			`
// push symbol THIS index 6 
@6
D=A
@THIS
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
`,
		},
		{
			"push that 5",
			args{
				cmd:     parser.PUSH,
				segment: "that",
				index:   5,
			},
			`
// push symbol THAT index 5 
@5
D=A
@THAT
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
`,
		},
		{
			"push pointer 0",
			args{
				cmd:     parser.PUSH,
				segment: "pointer",
				index:   0,
			},
			`
// push register R3
@R3
D=M
@SP
A=M
M=D
@SP
M=M+1
`,
		},
		{
			"push temp 6",
			args{
				cmd:     parser.PUSH,
				segment: "temp",
				index:   6,
			},
			`
// push register R11
@R11
D=M
@SP
A=M
M=D
@SP
M=M+1
`,
		},
		{
			"push static 8",
			args{
				cmd:     parser.PUSH,
				segment: "static",
				index:   8,
			},
			`
// push static Test.8
@Test.8
D=M
@SP
A=M
M=D
@SP
M=M+1
`,
		},
		{
			"pop local 0",
			args{
				cmd:     parser.POP,
				segment: "local",
				index:   0,
			},
			`
// pop symbol LCL index 0
@SP
M=M-1
@0
D=A
@LCL
D=D+M
@R13
M=D
@SP
A=M
D=M
@R13
A=M
M=D
`,
		},
		{
			"pop argument 2",
			args{
				cmd:     parser.POP,
				segment: "argument",
				index:   2,
			},
			`
// pop symbol ARG index 2
@SP
M=M-1
@2
D=A
@ARG
D=D+M
@R13
M=D
@SP
A=M
D=M
@R13
A=M
M=D
`,
		},
		{
			"pop this 6",
			args{
				cmd:     parser.POP,
				segment: "this",
				index:   6,
			},
			`
// pop symbol THIS index 6
@SP
M=M-1
@6
D=A
@THIS
D=D+M
@R13
M=D
@SP
A=M
D=M
@R13
A=M
M=D
`,
		},
		{
			"pop that 5",
			args{
				cmd:     parser.POP,
				segment: "that",
				index:   5,
			},
			`
// pop symbol THAT index 5
@SP
M=M-1
@5
D=A
@THAT
D=D+M
@R13
M=D
@SP
A=M
D=M
@R13
A=M
M=D
`,
		},
		{
			"pop temp 6",
			args{
				cmd:     parser.POP,
				segment: "temp",
				index:   6,
			},
			`
// pop register R11
@SP
M=M-1
@SP
A=M
D=M
@R11
M=D
`,
		},
		{
			"pop pointer 0",
			args{
				cmd:     parser.POP,
				segment: "pointer",
				index:   0,
			},
			`
// pop register R3
@SP
M=M-1
@SP
A=M
D=M
@R3
M=D
`,
		},
		{
			"pop static 8",
			args{
				cmd:     parser.POP,
				segment: "static",
				index:   8,
			},
			`
// pop static Test.8
@SP
M=M-1
@SP
A=M
D=M
@Test.8
M=D
`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := bytes.NewBufferString("")
			cw := &CodeWriter{
				w:  b,
				fn: "Test",
			}
			cw.WritePushPop(tt.args.cmd, tt.args.segment, tt.args.index)

			if string(b.Bytes()) != tt.want {
				t.Errorf("WritePushPop() = %s, want %v", b, tt.want)
			}
		})
	}
}

func TestCodeWriter_writePushConstant(t *testing.T) {
	tests := []struct {
		name string
		args int
		want string
	}{
		{
			"test",
			1,
			`
// push constant 1
@1
D=A
@SP
A=M
M=D
@SP
M=M+1
`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := bytes.NewBufferString("")
			cw := &CodeWriter{
				w: b,
			}
			cw.writePushConstant(tt.args)

			if string(b.Bytes()) != tt.want {
				t.Errorf("writePushConstant() = %s, want %v", b, tt.want)
			}
		})
	}
}

func TestCodeWriter_writePushSymbol(t *testing.T) {
	type args struct {
		symbol string
		index  int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"local",
			args{
				symbol: "LCL",
				index:  0,
			},
			`
// push symbol LCL index 0 
@0
D=A
@LCL
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
`,
		},
		{
			"argument",
			args{
				symbol: "ARG",
				index:  0,
			},
			`
// push symbol ARG index 0 
@0
D=A
@ARG
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
`,
		},
		{
			"this",
			args{
				symbol: "THIS",
				index:  0,
			},
			`
// push symbol THIS index 0 
@0
D=A
@THIS
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
`,
		},
		{
			"that",
			args{
				symbol: "THAT",
				index:  0,
			},
			`
// push symbol THAT index 0 
@0
D=A
@THAT
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
`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := bytes.NewBufferString("")
			cw := &CodeWriter{
				w: b,
			}
			cw.writePushSymbol(tt.args.symbol, tt.args.index)

			if string(b.Bytes()) != tt.want {
				t.Errorf("writePushSymbol() = %s, want %v", b, tt.want)
			}
		})
	}
}

func TestCodeWriter_writePopSymbol(t *testing.T) {
	type args struct {
		symbol string
		index  int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"local",
			args{
				symbol: "LCL",
				index:  0,
			},
			`
// pop symbol LCL index 0
@SP
M=M-1
@0
D=A
@LCL
D=D+M
@R13
M=D
@SP
A=M
D=M
@R13
A=M
M=D
`,
		},
		{
			"argument",
			args{
				symbol: "ARG",
				index:  0,
			},
			`
// pop symbol ARG index 0
@SP
M=M-1
@0
D=A
@ARG
D=D+M
@R13
M=D
@SP
A=M
D=M
@R13
A=M
M=D
`,
		},
		{
			"this",
			args{
				symbol: "THIS",
				index:  0,
			},
			`
// pop symbol THIS index 0
@SP
M=M-1
@0
D=A
@THIS
D=D+M
@R13
M=D
@SP
A=M
D=M
@R13
A=M
M=D
`,
		},
		{
			"that",
			args{
				symbol: "THAT",
				index:  0,
			},
			`
// pop symbol THAT index 0
@SP
M=M-1
@0
D=A
@THAT
D=D+M
@R13
M=D
@SP
A=M
D=M
@R13
A=M
M=D
`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := bytes.NewBufferString("")
			cw := &CodeWriter{
				w: b,
			}
			cw.writePopSymbol(tt.args.symbol, tt.args.index)

			if string(b.Bytes()) != tt.want {
				t.Errorf("writePopSymbol() = %s, want %v", b, tt.want)
			}
		})
	}
}

func TestCodeWriter_writePushRegister(t *testing.T) {
	tests := []struct {
		name  string
		index int
		want  string
	}{
		{
			"test",
			8,
			`
// push register R8
@R8
D=M
@SP
A=M
M=D
@SP
M=M+1
`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := bytes.NewBufferString("")
			cw := &CodeWriter{
				w: b,
			}
			cw.writePushRegister(tt.index)

			if string(b.Bytes()) != tt.want {
				t.Errorf("writePushRegister() = %s, want %v", b, tt.want)
			}
		})
	}
}

func TestCodeWriter_writePopRegister(t *testing.T) {
	tests := []struct {
		name  string
		index int
		want  string
	}{
		{
			"test",
			11,
			`
// pop register R11
@SP
M=M-1
@SP
A=M
D=M
@R11
M=D
`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := bytes.NewBufferString("")
			cw := &CodeWriter{
				w: b,
			}
			cw.writePopRegister(tt.index)

			if string(b.Bytes()) != tt.want {
				t.Errorf("writePopRegister() = %s, want %v", b, tt.want)
			}
		})
	}
}

func TestCodeWriter_writePushStatic(t *testing.T) {
	type args struct {
		index int
		fn    string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"push static 1",
			args{
				index: 1,
				fn:    "StackTest",
			},
			`
// push static StackTest.1
@StackTest.1
D=M
@SP
A=M
M=D
@SP
M=M+1
`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := bytes.NewBufferString("")
			cw := &CodeWriter{
				w:  b,
				fn: tt.args.fn,
			}
			cw.writePushStatic(tt.args.index)

			if string(b.Bytes()) != tt.want {
				t.Errorf("writePushStatic() = %s, want %v", b, tt.want)
			}
		})
	}
}

func TestCodeWriter_SetFileName(t *testing.T) {
	tests := []struct {
		name string
		fn   string
		want string
	}{
		{
			"test",
			"filename",
			"filename",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := bytes.NewBufferString("")
			cw := &CodeWriter{
				w: b,
			}
			cw.SetFileName(tt.fn)

			if !reflect.DeepEqual(cw.fn, tt.want) {
				t.Errorf("SetFileName() = %v, want %v", cw.fn, tt.want)
			}
		})
	}
}

func TestCodeWriter_writePopStatic(t *testing.T) {
	type args struct {
		index int
		fn    string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"pop static 0",
			args{
				index: 0,
				fn:    "StackTest",
			},
			`
// pop static StackTest.0
@SP
M=M-1
@SP
A=M
D=M
@StackTest.0
M=D
`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := bytes.NewBufferString("")
			cw := &CodeWriter{
				w:  b,
				fn: tt.args.fn,
			}
			cw.writePopStatic(tt.args.index)

			if string(b.Bytes()) != tt.want {
				t.Errorf("writePopStatic() = %s, want %v", b, tt.want)
			}
		})
	}
}

func TestCodeWriter_WriteLabel(t *testing.T) {
	type args struct {
		label string
		fn    string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"test1",
			args{
				label: "test",
				fn:    "TestFile",
			},
			`
// write label TestFile$test
(TestFile$test)
`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := bytes.NewBufferString("")
			cw := &CodeWriter{
				w:  b,
				fn: tt.args.fn,
			}
			cw.WriteLabel(tt.args.label)

			if string(b.Bytes()) != tt.want {
				t.Errorf("WriteLabel() = %s, want %v", b, tt.want)
			}
		})
	}
}

func TestCodeWriter_WriteIf(t *testing.T) {
	type args struct {
		label string
		fn    string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"test",
			args{
				label: "test",
				fn:    "TestFile",
			},
			`
// if-goto label TestFile$test
@SP
M=M-1
@SP
A=M
D=M
@TestFile$test
D;JNE
`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := bytes.NewBufferString("")
			cw := &CodeWriter{
				w:  b,
				fn: tt.args.fn,
			}
			cw.WriteIf(tt.args.label)

			if string(b.Bytes()) != tt.want {
				t.Errorf("WriteIf() = %s, want %v", b, tt.want)
			}
		})
	}
}

func TestCodeWriter_WriteGoto(t *testing.T) {
	type args struct {
		label string
		fn    string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"test",
			args{
				label: "test",
				fn:    "TestFile",
			},
			`
// goto label TestFile$test
@TestFile$test
0;JMP
`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := bytes.NewBufferString("")
			cw := &CodeWriter{
				w:  b,
				fn: tt.args.fn,
			}
			cw.WriteGoto(tt.args.label)

			if string(b.Bytes()) != tt.want {
				t.Errorf("WriteGoto() = %s, want %v", b, tt.want)
			}
		})
	}
}

func TestCodeWriter_WriteReturn(t *testing.T) {
	want := `
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
	t.Run("test1", func(t *testing.T) {
		b := bytes.NewBufferString("")
		cw := &CodeWriter{
			w:  b,
			fn: "test_fn",
		}
		cw.WriteReturn()

		if string(b.Bytes()) != want {
			t.Errorf("WriteReturn() = %s, want %v", b, want)
		}
	})
}

func TestCodeWriter_WriteFunction(t *testing.T) {
	type args struct {
		funcname string
		numlocal int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"test",
			args{
				funcname: "TestFunction",
				numlocal: 3,
			},
			`
// function TestFunction local nums 3
(TestFunction)

// push constant 0
@0
D=A
@SP
A=M
M=D
@SP
M=M+1

// push constant 0
@0
D=A
@SP
A=M
M=D
@SP
M=M+1

// push constant 0
@0
D=A
@SP
A=M
M=D
@SP
M=M+1
`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := bytes.NewBufferString("")
			cw := &CodeWriter{
				w: b,
			}
			cw.WriteFunction(tt.args.funcname, tt.args.numlocal)

			if string(b.Bytes()) != tt.want {
				t.Errorf("WriteFunction() = %s, want %v", b, tt.want)
			}
		})
	}
}
