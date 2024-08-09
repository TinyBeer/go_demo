package main

import (
	"fmt"
	"go/scanner"
	"go/token"
	"io"
	"os"
)

func main() {
	oFile, err := os.OpenFile("./hello/hello.go", os.O_RDONLY, 0777)
	if err != nil {
		panic(err)
	}
	src, err := io.ReadAll(oFile)
	if err != nil {
		panic(err)
	}
	fset := token.NewFileSet()
	file := fset.AddFile(oFile.Name(), fset.Base(), len(src))
	var s scanner.Scanner
	s.Init(file, src, nil, scanner.ScanComments)
	for {
		pos, tok, lit := s.Scan()
		fmt.Printf("%s\t%s\t%q\n", fset.Position(pos), tok, lit)
		if tok == token.EOF {
			break
		}
	}

}
