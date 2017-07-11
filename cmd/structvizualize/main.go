package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"

	"github.com/ilijamt/structvizualizer"
)

func main() {

	fset := token.NewFileSet() // positions are relative to fset
	files := make(map[string]*ast.File)

	for _, arg := range os.Args[1:] {
		f, err := parser.ParseFile(fset, arg, nil, parser.AllErrors)
		if err != nil {
			fmt.Println(err)
			return
		}
		files[arg] = f
	}

	assembler := structvizualizer.NewAssembler("G")

	for filename, f := range files {
		assembler.Add(filename, f)
	}

	assembler.Parse()
	assembler.Template()
	//assembler.Dump()

}
