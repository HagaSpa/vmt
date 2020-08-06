package parser

import (
	"reflect"
	"strings"
	"testing"
)

func TestParser_HasMoreCommands(t *testing.T) {
	tests := []struct {
		name string
		args string
		want bool
	}{
		{
			"true",
			"abcdefg",
			true,
		},
		{
			"false",
			"",
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := strings.NewReader(tt.args)
			p := New(b)

			if got := p.HasMoreCommands(); got != tt.want {
				t.Errorf("Parser.HasMoreCommands() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParser_Advance(t *testing.T) {
	tests := []struct {
		name string
		args string
		want string
	}{
		{
			"test1",
			"test",
			"test",
		},
		{
			"comment",
			"//comment",
			"",
		},
		{
			"whitespace",
			"  test  ",
			"test",
		},
		{
			"empty",
			"",
			"",
		},
		{
			"comment_in_sentence",
			"test // comment",
			"test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := strings.NewReader(tt.args)
			p := New(b)
			p.HasMoreCommands()
			p.Advance()

			if !reflect.DeepEqual(p.input, tt.want) {
				t.Errorf("Parser.Advance() = %s, want %s", p.input, tt.want)
			}
		})
	}
}

func TestParser_CommandType(t *testing.T) {
	tests := []struct {
		name string
		args string
		want Type
	}{
		{
			"add",
			"add",
			ARITHMETIC,
		},
		{
			"sub",
			"sub",
			ARITHMETIC,
		},
		{
			"neg",
			"neg",
			ARITHMETIC,
		},
		{
			"eq",
			"eq",
			ARITHMETIC,
		},
		{
			"gt",
			"gt",
			ARITHMETIC,
		},
		{
			"lt",
			"lt",
			ARITHMETIC,
		},
		{
			"and",
			"and",
			ARITHMETIC,
		},
		{
			"or",
			"or",
			ARITHMETIC,
		},
		{
			"not",
			"not",
			ARITHMETIC,
		},
		{
			"push",
			"push",
			PUSH,
		},
		{
			"pop",
			"pop",
			POP,
		},
		{
			"label",
			"label",
			LABEL,
		},
		{
			"goto",
			"goto",
			GOTO,
		},
		{
			"if-goto",
			"if-goto",
			IF,
		},
		{
			"function",
			"function",
			FUNCTION,
		},
		{
			"return",
			"return",
			RETURN,
		},
		{
			"call",
			"call",
			CALL,
		},
		{
			"invalid vm command",
			"test",
			None,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := strings.NewReader(tt.args)
			p := New(b)
			p.HasMoreCommands()
			p.Advance()

			if got := p.CommandType(); got != tt.want {
				t.Errorf("Parser.CommandType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParser_Arg1(t *testing.T) {
	tests := []struct {
		name string
		args string
		want string
	}{
		// ARITHMETIC
		{
			"add",
			"add",
			"add",
		},
		{
			"sub",
			"sub",
			"sub",
		},
		{
			"neg",
			"neg",
			"neg",
		},
		{
			"eq",
			"eq",
			"eq",
		},
		{
			"gt",
			"gt",
			"gt",
		},
		{
			"lt",
			"lt",
			"lt",
		},
		{
			"and",
			"and",
			"and",
		},
		{
			"or",
			"or",
			"or",
		},
		{
			"not",
			"not",
			"not",
		},
		// PUSH
		{
			"push",
			"push constant 7",
			"constant",
		},
		// POP
		{
			"pop",
			"pop local 0",
			"local",
		},
		// LABEL
		{
			"label",
			"label LOOP",
			"LOOP",
		},
		// GOTO
		{
			"goto",
			"goto LOOP",
			"LOOP",
		},
		// IF
		{
			"if",
			"if-goto END",
			"END",
		},
		// FUNCTION
		{
			"function",
			"function mult 2",
			"mult",
		},

		/*
			when RETURN, Args1() must not be called...
		*/

		// CALL
		{
			"call",
			"call mult 2",
			"mult",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := strings.NewReader(tt.args)
			p := New(b)
			p.HasMoreCommands()
			p.Advance()

			if got := p.Arg1(); got != tt.want {
				t.Errorf("Parser.Arg1() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParser_Arg2(t *testing.T) {
	tests := []struct {
		name    string
		args    string
		want    int
		wantErr bool
	}{
		// PUSH
		{
			"push",
			"push constant 7",
			7,
			false,
		},
		// POP
		{
			"pop",
			"pop local 0",
			0,
			false,
		},
		// FUNCTION
		{
			"function",
			"function mult 2",
			2,
			false,
		},
		// CALL
		{
			"call",
			"call mult 2",
			2,
			false,
		},

		/*
			when ARITHMETIC, LABEL, GOTO, IF and RETURN, Args2() must not be called...
		*/

		// invalid arg2
		{
			"invalid arg2",
			"push constant test",
			0,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := strings.NewReader(tt.args)
			p := New(b)
			p.HasMoreCommands()
			p.Advance()

			got, err := p.Arg2()
			if (err != nil) != tt.wantErr {
				t.Errorf("Parser.Arg2() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Parser.Arg2() = %v, want %v", got, tt.want)
			}
		})
	}
}
