package codewriter

import (
	"bytes"
	"strings"
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
			p := parser.New(strings.NewReader(tt.line))
			cw := &CodeWriter{
				w: b,
				p: p,
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
			p := parser.New(strings.NewReader(tt.line))
			cw := &CodeWriter{
				w: b,
				p: p,
			}
			cw.writeUnaryOperator(tt.args)

			if string(b.Bytes()) != tt.want {
				t.Errorf("writeUnaryOperator() = %s, want %v", b, tt.want)
			}
		})
	}
}

func TestCodeWriter_writeConditionOperator(t *testing.T) {
	type args struct {
		op string
		n  int
	}
	tests := []struct {
		name string
		line string
		args args
		want string
	}{
		{
			"eq",
			"eq",
			args{
				op: "eq",
				n:  1,
			},
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
			args{
				op: "gt",
				n:  1,
			},
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
			args{
				op: "lt",
				n:  1,
			},
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
			p := parser.New(strings.NewReader(tt.line))
			cw := &CodeWriter{
				w: b,
				p: p,
			}
			cw.writeConditionOperator(tt.args.op, tt.args.n)

			if string(b.Bytes()) != tt.want {
				t.Errorf("writeConditionOperator() = %s, want %v", b, tt.want)
			}
		})
	}
}
