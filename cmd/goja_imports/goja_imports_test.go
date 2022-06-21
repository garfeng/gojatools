package main

import (
	"fmt"
	"os"
	"testing"
)

func TestImport(t *testing.T) {
	os.Mkdir("hello", 0777)
	err := importGoPackage("image", "hello", "hello/hello.go", "hello/hello.ts")
	if err != nil {
		fmt.Println(err)
	}
}
