package main

import (
	"flag"
	"fmt"
	"os"
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

	// generate parser
	p := parser.New(f)
	for p.HasMoreCommands() {
		p.Advance()
		fmt.Println(p)
		// TODO: switch Type
	}

}
