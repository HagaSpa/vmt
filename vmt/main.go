package main

import (
	"flag"
	"os"
	"path/filepath"
	"regexp"
	"vmt/codewriter"
	"vmt/parser"
)

func main() {
	// parse args
	flag.Parse()
	flags := flag.Args()
	if flags == nil {
		os.Exit(1)
	}

	// generate asm
	rep := regexp.MustCompile(`.vm$`)
	name := filepath.Base(rep.ReplaceAllString(flags[0], "")) + ".asm"
	asm, err := os.Create(name)
	if err != nil {
		os.Exit(1)
	}
	defer asm.Close()
	cw := codewriter.New(asm)

	// open vm
	// TODO: corresponds multiple vm file in directory
	f, err := os.Open(flags[0])
	if err != nil {
		os.Exit(1)
	}
	defer f.Close()
	p := parser.New(f)
	cw.SetFileName(rep.ReplaceAllString(flags[0], ""))

	// write assembley
	for p.HasMoreCommands() {
		p.Advance()
		switch p.CommandType() {
		case parser.ARITHMETIC:
			cw.WriteArithmetic(p.Arg1())
		case parser.PUSH, parser.POP:
			index, err := p.Arg2()
			if err != nil {
				os.Exit(1)
			}
			cw.WritePushPop(p.CommandType(), p.Arg1(), index)
		}
	}

}
