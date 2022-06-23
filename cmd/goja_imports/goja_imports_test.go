package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/dop251/goja"
)

func TestImport(t *testing.T) {
	os.RemoveAll("hello")
	os.Mkdir("hello", ModDir)
	err := importGoPackage("image", "hello", "hello/hello.go", "hello/hello.ts")
	if err != nil {
		fmt.Println(err)
	}
}

type MyError struct{}

func (m *MyError) Error() string {
	return "hello"
}

func fnDemo(in1, in2 int, in3 ...float64) (int, float32, error) {
	b := in1 + in2
	if b == 10 {
		return b, 1, new(MyError)
	}
	return b, 1, nil
}

func fnDemo2() error {
	return nil
}

func TestGojaErr(t *testing.T) {
	engine := goja.New()
	engine.Set("fnDemo", fnDemo)
	engine.Set("fnDemo2", fnDemo2)
	engine.Set("console", map[string]interface{}{
		"log": fmt.Println,
	})

	script := `
	try{
	var c = fnDemo2();
	console.log("js:", c);
	var b = fnDemo(1,5,2)
	console.log("js:",  b)
	} catch(e) {
		console.log("jsErr:", e);
	}
	`

	res, err := engine.RunString(script)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("res:", res)
}
