package parser

import (
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
