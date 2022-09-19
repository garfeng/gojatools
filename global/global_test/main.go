package main

import (
	"fmt"
	"io/ioutil"

	"github.com/dop251/goja"
	_ "github.com/garfeng/gojatools/global/js"
	"github.com/garfeng/gojatools/importgo"
)

func main() {
	engine := goja.New()
	importgo.ExportConsoleToEngine(engine)
	importgo.ExportRequireToEngine(engine)
	importgo.ExportToEngine(engine)
	buff, _ := ioutil.ReadFile("./test.js")
	code := string(buff)
	code = importgo.FormatCode(code)

	v, err := engine.RunString(code)
	engine.RunString(code)
	engine.RunString(code)
	engine.RunString(code)

	if err != nil {
		fmt.Println(err)
		return
	}

	if v != goja.Undefined() {
		fmt.Println(v)
	}

}
