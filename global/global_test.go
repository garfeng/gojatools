package global

import (
	"fmt"
	"testing"

	"github.com/dop251/goja"
)

func TestGlobalValues(t *testing.T) {
	engine := goja.New()
	engine.Set("Load", Load)
	engine.Set("Save", Save)
	engine.Set("log", fmt.Println)

	script := `
	// initial value
var data = {
    name: "Alice",
    id: 1,
    age: 10,
};

// Load from go engine if existed
data = Load("alice") || data;

// update
data.age ++

Save("alice", data)

log(data)
`

	engine.RunScript("script", script)
	engine.RunScript("script", script)
	engine.RunScript("script", script)
}
