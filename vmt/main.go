package main

import (
	"flag"
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
		log.Fatalln(err.Error())
	}
	defer asm.Close()

	// generate codewriter
	cw := codewriter.New(asm)

	// IsDir?
	if !fInfo.IsDir() {
		cw.SetFileName(bname)
		translate(flags[0], cw)
		log.Println("translated vm: " + bname + ".asm")
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
	log.Println("translated multiple vm: " + bname + ".asm")
}

func translate(vmn string, cw *codewriter.CodeWriter) {
	// open vm
	f, err := os.Open(vmn)
	if err != nil {
		log.Fatalln(err.Error())
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
				log.Fatalln(err.Error())
			}
			cw.WritePushPop(p.CommandType(), p.Arg1(), index)
		case parser.LABEL:
			cw.WriteLabel(p.Arg1())
		case parser.IF:
			cw.WriteIf(p.Arg1())
		case parser.GOTO:
			cw.WriteGoto(p.Arg1())
		case parser.FUNCTION:
			numlocal, err := p.Arg2()
			if err != nil {
				log.Fatalln(err.Error())
			}
			cw.WriteFunction(p.Arg1(), numlocal)
		case parser.RETURN:
			cw.WriteReturn()
		}
	}
}
