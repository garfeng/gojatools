package importgo

import (
	"fmt"
	"io/ioutil"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"

	"github.com/dop251/goja"
)

var (
	pkgs = map[string]Package{}
)

type Package map[string]interface{}

func Add(name string, namespace string, pkgCore Package) {
	pkgs[name] = pkgCore
}

func Remove(name string) {
	delete(pkgs, name)
}

func Import(name string) Package {
	return pkgs[name]
}

func Pkgs() map[string]Package {
	return pkgs
}

type requireEngine struct {
	engine  *goja.Runtime
	modules map[string]goja.Value
}

func (r *requireEngine) Require(name string) interface{} {
	_, jsPkgName := filepath.Split(name)
	jsPkgName = strings.ToLower(jsPkgName)
	jsExt := filepath.Ext(jsPkgName)

	if jsExt == ".js" || jsExt == ".ts" {
		jsPkgName = jsPkgName[:len(jsPkgName)-len(jsExt)]
	}

	if v, find := pkgs[jsPkgName]; find {
		return v
	}

	if v, find := r.modules[name]; find {
		return v
	}

	buff, err := ioutil.ReadFile(name)
	if err != nil {
		panic("package not find: " + name)
	}
	code := string(buff)

	if strings.ToLower(path.Ext(name)) == ".json" {
		code = "module.exports = JSON.parse('" + template.JSEscapeString(code) + "')"
	}

	source := "(function() {\r\n let module = {exports:{}};\r\n" + code + "\n return module.exports;\r\n })()"

	source = FormatCode(source)

	v, err := r.engine.RunString(source)

	if err != nil {
		fmt.Println(err)
	}

	return v.Export()
}

func ExportRequireToEngine(engine *goja.Runtime) {
	re := &requireEngine{
		engine: engine,
	}
	engine.Set("require", re.Require)
}

func ExportConsoleToEngine(engine *goja.Runtime) {
	engine.Set("console", Package{
		"log": fmt.Println,
	})
}

func ExportToEngine(engine *goja.Runtime) {
	for k, v := range pkgs {
		engine.Set(k, v)
	}
}

var (
	importJSRegexp = regexp.MustCompile(`import (.+) from (.+)`)
)

func FormatCode(code string) string {
	//	matched := importJSRegexp.FindAllStringSubmatch(code, -1)
	code = importJSRegexp.ReplaceAllStringFunc(code, func(s string) string {
		s = strings.TrimSpace(s)
		s = strings.Trim(s, ";")
		s = strings.Replace(s, "import", "const", -1)
		s = strings.Replace(s, "from", "= require(", -1)
		s += ");"

		return s
	})
	return code
}
