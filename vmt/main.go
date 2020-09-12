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
	if len(flags) == 0 {
		log.Fatalln("Please specify the command argument for vm name")
	}

	// invalid args?
	fInfo, err := os.Stat(flags[0])
	if err != nil {
		log.Fatalln(err.Error())
	}

	// generate asm
	rep := regexp.MustCompile(`.vm$`)
	bname := filepath.Base(rep.ReplaceAllString(flags[0], ""))
	asm, err := os.Create(bname + ".asm")
	if err != nil {
		os.Exit(1)
	}
	defer asm.Close()

	// generate codewriter
	cw := codewriter.New(asm)

	// IsDir?
	if !fInfo.IsDir() {
		cw.SetFileName(bname)
		translate(flags[0], cw)
		return
	}

	// multiple files in directory
	fpath := flags[0] + "/*.vm"
	files, err := filepath.Glob(fpath)
	for _, f := range files {
		pf := filepath.Base(rep.ReplaceAllString(f, ""))
		cw.SetFileName(pf)
		translate(f, cw)
	}
	fmt.Println("IsDir")
}

func translate(vmn string, cw *codewriter.CodeWriter) {
	// open vm
	f, err := os.Open(vmn)
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
