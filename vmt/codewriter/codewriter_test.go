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
		line string
		op   string
		want string
	}{
		{
			"eq",
			"eq",
			"eq",
			`
// Condition Operator eq
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
@1
D;JEQ
@SP
A=M
M=0
(1)
@SP
M=M+1
`,
		},
		{
			"gt",
			"gt",
			"gt",
			`
// Condition Operator gt
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
@1
D;JGT
@SP
A=M
M=0
(1)
@SP
M=M+1
`,
		},
		{
			"lt",
			"lt",
			"lt",
			`
// Condition Operator lt
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
@1
D;JLT
@SP
A=M
M=0
(1)
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
		line string
		op   string
		want string
	}{
		{
			"eq",
			"eq",
			"eq",
			`
// Condition Operator eq
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
@1
D;JEQ
@SP
A=M
M=0
(1)
@SP
M=M+1
`,
		},
		{
			"gt",
			"gt",
			"gt",
			`
// Condition Operator eq
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
@1
D;JEQ
@SP
A=M
M=0
(1)
@SP
M=M+1

// Condition Operator gt
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
@2
D;JGT
@SP
A=M
M=0
(2)
@SP
M=M+1
`,
		},
		{
			"lt",
			"lt",
			"lt",
			`
// Condition Operator eq
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
@1
D;JEQ
@SP
A=M
M=0
(1)
@SP
M=M+1

// Condition Operator gt
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
@2
D;JGT
@SP
A=M
M=0
(2)
@SP
M=M+1

// Condition Operator lt
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
@3
D;JLT
@SP
A=M
M=0
(3)
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
		{
			"eq",
			"eq",
			`
// Condition Operator eq
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
@1
D;JEQ
@SP
A=M
M=0
(1)
@SP
M=M+1
`,
		},
		{
			"gt",
			"gt",
			`
// Condition Operator gt
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
@1
D;JGT
@SP
A=M
M=0
(1)
@SP
M=M+1
`,
		},
		{
			"lt",
			"lt",
			`
// Condition Operator lt
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
@1
D;JLT
@SP
A=M
M=0
(1)
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
