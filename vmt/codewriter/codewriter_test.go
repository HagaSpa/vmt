package codewriter

import (
	"bytes"
	"strings"
	"testing"
	"vmt/parser"
)

func TestCodeWriter_writeBinaryOperator(t *testing.T) {

	tests := []struct {
		name    string
		command string
		args    string
		want    string
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
			p := parser.New(strings.NewReader(tt.command))
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
