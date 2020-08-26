package codewriter

import (
	"bytes"
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := bytes.NewBufferString("")
			cw := &CodeWriter{
				w: b,
			}
			cw.WritePushPop(tt.args.cmd, tt.args.segment, tt.args.index)
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
