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
