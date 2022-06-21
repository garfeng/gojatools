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
)

func init() {
	flag.Parse()

}

func main() {
	fmt.Println(os.Args[0], "--help for help")
	if pkg == nil || (*pkg) == "" {
		flag.PrintDefaults()
		return
	}
	fmt.Println("Current Args:")
	fmt.Println("-p", *pkg, "-i", *inPkg, "-go", *goPath, "-ja", *jsPath)

	base, _ := filepath.Split(*goPath)
	os.MkdirAll(base, 0777)
	base, _ = filepath.Split(*jsPath)
	os.MkdirAll(base, 0777)

	err := importGoPackage(*inPkg, *pkg, *goPath, *jsPath)
	if err != nil {
		fmt.Println(err)
	}
}
