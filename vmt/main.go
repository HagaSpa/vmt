package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"vmt/parser"
)

func main() {
	// parse args
	flag.Parse()
	flags := flag.Args()
	if flags == nil {
		os.Exit(1)
	}

	// open vm
	f, err := os.Open(flags[0])
	if err != nil {
		os.Exit(1)
	}
	defer f.Close()

	// generate asm
	rep := regexp.MustCompile(`.vm$`)
	name := filepath.Base(rep.ReplaceAllString(flags[0], "")) + ".asm"
	asm, err := os.Create(name)
	if err != nil {
		os.Exit(1)
	}
	defer asm.Close()

	// generate parser
	p := parser.New(f)
	for p.HasMoreCommands() {
		p.Advance()
		fmt.Println(p)
		// TODO: switch Type
	}

}
