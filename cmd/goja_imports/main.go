package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

var (
	pkg    = flag.String("p", "", "dst package name")
	inPkg  = flag.String("i", "", "input package path")
	goPath = flag.String("go", "js/js.go", "dst go file")
	jsPath = flag.String("ja", "js/js.ts", "dst typescript file")
	bigInt = flag.Bool("bigInt", false, "has big int")
)

func main() {
	flag.Parse()

	if pkg == nil || (*pkg) == "" {
		fmt.Println(os.Args[0], "--help for help")
		flag.PrintDefaults()
		return
	}

	base, _ := filepath.Split(*goPath)
	os.MkdirAll(base, ModDir)
	base, _ = filepath.Split(*jsPath)
	os.MkdirAll(base, ModDir)

	err := importGoPackage(*inPkg, *pkg, *goPath, *jsPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
