package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/dop251/goja"
	_ "github.com/garfeng/gojatools/examples/gocv/js"
	"github.com/garfeng/gojatools/importgo"
	"gocv.io/x/gocv"
)

func GetAddr(v gocv.Mat) *gocv.Mat {
	return &v
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage:", os.Args[0], "xxx.js")
		return
	}
	engine := goja.New()
	importgo.ExportConsoleToEngine(engine)
	importgo.ExportRequireToEngine(engine)
	importgo.ExportToEngine(engine)

	engine.Set("addr", GetAddr)

	buff, _ := ioutil.ReadFile(os.Args[1])
	code := string(buff)
	code = importgo.FormatCode(code)
	v, err := engine.RunScript(os.Args[1], code)
	if err != nil {
		fmt.Println(err)
		return
	}
	if v != goja.Undefined() {
		fmt.Println(v)
	}
}
