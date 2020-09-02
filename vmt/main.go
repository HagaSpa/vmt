package main

import (
	"flag"
	"fmt"
	"log"
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
	bname := filepath.Base(rep.ReplaceAllString(flags[0], ""))
	name := bname + ".asm"
	asm, err := os.Create(name)
	if err != nil {
		os.Exit(1)
	}
	defer asm.Close()
	cw := codewriter.New(asm)

	// IsDir?
	fInfo, err := os.Stat(flags[0])
	if err != nil {
		log.Fatalln(err)
	}
	if fInfo.IsDir() {
		// TODO: corresponds multiple files in directory
		fmt.Println("IsDir")
		return
	}

	cw.SetFileName(bname)
	translate(flags[0], cw)
}

func translate(fn string, cw *codewriter.CodeWriter) {
	// open vm
	f, err := os.Open(fn)
	if err != nil {
		os.Exit(1)
	}
	defer f.Close()

	// write assembley
	p := parser.New(f)
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
