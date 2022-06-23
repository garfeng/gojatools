package main

import (
	"bytes"
	"embed"
	"fmt"
	"go/ast"
	"go/format"
	"io/ioutil"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/wzshiming/gotype"
)

type TypeDefine struct {
	Name    string
	Doc     string
	Comment string

	Fields []ValueDefine
	Funcs  []FuncDefine

	Type gotype.Type
}

type ValueDefine struct {
	Value          string
	Name           string
	Doc            string
	Comment        string
	TypeName       string
	Type           gotype.Type
	ShouldCastType bool
	SrcPkg         string
	IsGlobal       bool
}

type NameAndType struct {
	Name     string
	TypeName string
	Type     gotype.Type
}

type FuncDefine struct {
	Name    string
	Doc     string
	Comment string
	Type    gotype.Type

	Ins  []NameAndType
	Outs []NameAndType
}

func importGoPackage(pkgPath, dstPkgName, dstGoFile, dstJSFile string) error {
	os.Mkdir(tmpPath, ModDir)
	defer os.RemoveAll(tmpPath)

	imp := gotype.NewImporter()
	pkgInfos, err := imp.Import(pkgPath, "")
	if err != nil {
		return err
	}

	values := []ValueDefine{}
	types := []TypeDefine{}
	funcs := []FuncDefine{}

	num := pkgInfos.NumChild()

	pkgPath = strings.Trim(pkgPath, "/")

	globalImportPkgs := map[string]bool{}

	for i := 0; i < num; i++ {
		v := pkgInfos.Child(i)
		if ast.IsExported(v.Name()) {
			switch v.Kind() {
			case gotype.Declaration:
				decl := v.Declaration()
				switch decl.Kind() {
				case gotype.Func:
					ImportFunc(v, globalImportPkgs, &funcs)
				default:
					ImportValues(v, pkgPath, true, &values)
				}
			case gotype.Scope, gotype.Invalid:
			default:
				ImportTypes(v, pkgPath, globalImportPkgs, &types)
			}
		}
	}

	globalImportPkgsList := []string{}

	for k := range globalImportPkgs {
		if k != pkgPath {
			globalImportPkgsList = append(globalImportPkgsList, k)
		}
	}

	buf := bytes.NewBuffer(nil)

	_, exportName := filepath.Split(pkgPath)
	exportNameTitle := strings.ToUpper(pkgPath)
	exportNameTitle = strings.ReplaceAll(exportNameTitle, ".", "_")
	exportNameTitle = strings.ReplaceAll(exportNameTitle, "/", "_")

	_, jsPkgName := filepath.Split(dstJSFile)
	jsPkgName = strings.ToLower(jsPkgName)
	jsExt := filepath.Ext(jsPkgName)

	jsPkgName = jsPkgName[:len(jsPkgName)-len(jsExt)]

	_, jsNamespace := filepath.Split(pkgPath)

	_, jsFilename := filepath.Split(dstJSFile)
	_, goFilename := filepath.Split(dstGoFile)

	tmplData := map[string]interface{}{
		"SrcPkgName":      pkgPath,
		"DstPkgName":      dstPkgName,
		"GoFile":          dstGoFile,
		"JSFile":          dstJSFile,
		"JSPkgName":       jsPkgName,
		"JSNamespace":     jsNamespace,
		"ExportName":      exportName,
		"ExportNameTitle": exportNameTitle,
		"TypeSpec":        types,
		"ValueSpec":       values,
		"FuncSpec":        funcs,
		"JSFilename":      jsFilename,
		"GoFilename":      goFilename,
		"PkgImports":      globalImportPkgsList,
	}

	templateFile.ExecuteTemplate(buf, "tempGo", tmplData)

	code := buf.String()

	codeFormated, err := format.Source([]byte(code))

	if err != nil {
		ioutil.WriteFile(dstGoFile, []byte(code), ModeCode)
		fmt.Println("Err of format code: ", err)
		return err
	}

	ioutil.WriteFile(dstGoFile, codeFormated, ModeCode)

	bufJs := bytes.NewBuffer(nil)
	jsTemplate.ExecuteTemplate(bufJs, "tempJS", tmplData)
	ioutil.WriteFile(dstJSFile, bufJs.Bytes(), ModeCode)

	return nil
}

var (
	tmpPath = "./.tmp_" + fmt.Sprint(time.Now().UnixNano())
)

func ImportFunc(v gotype.Type, globalImportPkgs map[string]bool, fnList *[]FuncDefine) {

	fn := FuncDefine{
		Name:    v.Name(),
		Doc:     comment(v.Doc().Text()),
		Comment: comment(v.Comment().Text()),
		Type:    v,
	}
	decl := v.Declaration()
	if decl.NumIn() > 0 {
		fn.Ins = []NameAndType{}
		for i := 0; i < decl.NumIn(); i++ {
			in := decl.In(i)
			fn.Ins = append(fn.Ins, NameAndType{
				Name:     in.Name(),
				TypeName: in.Declaration().String(),
				Type:     in.Declaration(),
			})
		}
	}

	if decl.NumOut() > 0 {
		fn.Outs = []NameAndType{}
		for i := 0; i < decl.NumOut(); i++ {
			out := decl.Out(i)
			fn.Outs = append(fn.Outs, NameAndType{
				Name:     out.Name(),
				TypeName: out.Declaration().String(),
				Type:     out.Declaration(),
			})
		}
	}

	GoFuncImportsPkgs(fn, globalImportPkgs)

	(*fnList) = append(*fnList, fn)
}

func ImportTypes(v gotype.Type, pkgPath string, globalImportPkgs map[string]bool, typeList *[]TypeDefine) {
	tp := TypeDefine{
		Name:    v.Name(),
		Doc:     comment(v.Doc().Text()),
		Comment: comment(v.Comment().Text()),
		Type:    v,
	}

	if v.Kind() == gotype.Struct {
		if v.NumField() > 0 {
			fields := []ValueDefine{}
			for i := 0; i < v.NumField(); i++ {
				para := v.Field(i)
				if ast.IsExported(para.Name()) {
					ImportValues(para, pkgPath, false, &fields)
				}
			}
			tp.Fields = fields
		}
	}

	methods := []FuncDefine{}
	ImportSubInterafaceMethods(v, globalImportPkgs, &methods, 0)

	tp.Funcs = methods

	(*typeList) = append((*typeList), tp)
}

func ImportSubInterafaceMethods(v gotype.Type, globalPkg map[string]bool, methods *[]FuncDefine, padding int) {
	if v.Kind() != gotype.Interface {
		if v.NumMethod() > 0 {
			for i := 0; i < v.NumMethod(); i++ {
				fn := v.Method(i)
				if fn.Kind() == gotype.Declaration {
					if ast.IsExported(fn.Name()) {
						ImportFunc(fn, globalPkg, methods)
					}
				} else {
					//	fmt.Println("not decl", fn.Kind())
				}

			}
		}
	}

	if v.Kind() == gotype.Struct {
		num := v.NumField()
		for i := 0; i < num; i++ {
			one := v.Field(i)
			if one.IsAnonymous() {
				ImportSubInterafaceMethods(one, globalPkg, methods, padding+2)
			}
		}
	}

	prefix := ""

	for i := 0; i < padding; i++ {
		prefix += "  "
	}

	if v.Kind() == gotype.Interface {
		//	fmt.Println(prefix, v, v.NumField(), "===>")

		num := v.NumField()
		for i := 0; i < num; i++ {
			one := v.Field(i)
			if one.Kind() == gotype.Declaration {
				decl := one.Declaration()
				//	fmt.Println(prefix, "  ", "Field:", one.Name(), decl.NumMethod(), decl.Kind())

				if decl.Kind() == gotype.Interface {
					ImportSubInterafaceMethods(decl, globalPkg, methods, padding+2)
				} else {
					//	fn := decl
					//			fmt.Println(prefix, "  ", "Import["+one.Name()+"]", fn.Kind(), ast.IsExported(one.Name()))
					if ast.IsExported(one.Name()) {
						ImportFunc(one, globalPkg, methods)
					}
				}
			} else {
				//ImportSubInterafaceMethods(one, globalPkg, methods, padding+2)
			}
		}
	}

}

func GetRealDecl(v gotype.Type) gotype.Type {
	if v.Kind() == gotype.Declaration {
		return GetRealDecl(v.Declaration())
	}
	return v
}

func ImportValues(v gotype.Type, pkgPath string, isGlobal bool, valueList *[]ValueDefine) {
	typeName := ""

	if v.Kind() == gotype.Declaration {
		typeName = GetRealDecl(v).Kind().String()
	} else {
		name := v.Name()
		s := v.String()
		typeName = strings.TrimSpace(strings.Replace(s, name, "", 1))
	}

	tp := ValueDefine{
		Value:    v.Value(),
		Name:     v.Name(),
		Doc:      comment(v.Doc().Text()),
		Comment:  comment(v.Comment().Text()),
		Type:     v,
		TypeName: typeName,
		SrcPkg:   pkgPath,
		IsGlobal: isGlobal,
	}

	if isGlobal {
		if (strings.ToLower(tp.TypeName) == "int") && tp.Value != "" {
			//		fmt.Println(tp.TypeName, tp.Name, "=", tp.Value)
			value, _ := strconv.ParseFloat(tp.Value, 64)
			if value > float64(math.MaxInt) || value < float64(math.MinInt) {
				tp.ShouldCastType = true
			}
		}

		if (strings.ToLower(tp.TypeName) == "int") && tp.Value == "" {
			tp.ShouldCastType = true
		}
	}

	(*valueList) = append((*valueList), tp)
}

func comment(s string) string {
	ss := []string{}
	dd := strings.Split(s, "\n")
	for i, v := range dd {
		if v == "" && i == len(dd)-1 {
			continue
		}
		v = strings.TrimSpace(v)
		if strings.Index(v, "//") != 0 {
			v = "// " + v
		}
		ss = append(ss, v)
	}

	return strings.Join(ss, "\n")
}

var templateFile = template.New("tempGo")
var jsTemplate = template.New("tempJS")

func initTemplate() {
	templateFile.Funcs(template.FuncMap{
		"GoNumFormat":  GoNumFormat,
		"GoFuncFormat": GoFuncFormat,
	})
	templateFile.Parse(goTemp)

	jsTemplate.Funcs(template.FuncMap{
		"jsInsFormat":    JsInsFormat,
		"jsOutsFormat":   JsOutsFormat,
		"jsTypeName":     jsTypeName,
		"JsImportFormat": JsImportFormat,
		"JsImportName":   JsImportName,
	})
	jsTemplate.Parse(jsTemp)
}

//go:embed tmpl
var goCodeTmpl embed.FS

var (
	goTemp, jsTemp string
)

func init() {
	var err error
	var buff []byte
	buff, err = goCodeTmpl.ReadFile("tmpl/gocode.tmpl")
	if err != nil {
		panic(err)
	}
	goTemp = string(buff)

	buff, err = goCodeTmpl.ReadFile("tmpl/jscode.tmpl")
	if err != nil {
		panic(err)
	}
	jsTemp = string(buff)

	initTemplate()
}

func GoNumFormat(v ValueDefine) string {
	//_, pkgName := filepath.Split(v.SrcPkg)

	if !v.ShouldCastType || !v.IsGlobal {
		return "origin." + v.Name
	}

	codeTmpl := `
		package main
		import (
			origin "%s"
			"fmt"
		)
		func main(){
			fmt.Println(origin.%s)
		}
	`
	code := fmt.Sprintf(codeTmpl, v.SrcPkg, v.Name)

	stdout, _, err := runGoCode(code)

	if err != nil {
		codeTmpl2 := `
		package main
		import (
			origin "%s"
			"fmt"
		)
		func main() {
			if origin.%s >= 0 {
				fmt.Print(1)
			} else {
				fmt.Print(-1)
			}
		}
		`
		code = fmt.Sprintf(codeTmpl2, v.SrcPkg, v.Name)
		stdout, stderr, err := runGoCode(code)
		if err != nil {
			fmt.Println(stderr)
			fmt.Println(stdout)
			fmt.Println(err)
			panic(err)
		}
		vfloat, _ := strconv.ParseFloat(stdout, 64)
		if vfloat >= 0 {
			fmt.Println("parse uint64", v.SrcPkg, v.Name, stdout)
			return fmt.Sprintf("uint64(origin.%s)", v.Name)
		} else {
			fmt.Println("parse int64", v.SrcPkg, v.Name, stdout)
			return fmt.Sprintf("int64(origin.%s)", v.Name)
		}
	}

	fmt.Println("parse direct", v.SrcPkg, v.Name, stdout)
	return "origin." + v.Name
}

func runGoCode(code string) (stdout, stderr string, err error) {

	codeName := filepath.Join(tmpPath, "main.go")
	ioutil.WriteFile(codeName, []byte(code), ModeCode)

	cmd := exec.Command("go", "run", codeName)
	wOut := bytes.NewBuffer(nil)
	wErr := bytes.NewBuffer(nil)

	cmd.Stdout = wOut
	cmd.Stderr = wErr
	err = cmd.Run()

	return wOut.String(), wErr.String(), err
}

func GoNumFormat2(v ValueDefine) string {
	//return "origin." + v.Name

	name := "origin." + v.Name
	if !v.ShouldCastType {
		return name
	}

	const div int64 = 0x100000000
	// const and = 0xFFFFFFFF

	// value - (value / (div)) * div
	return fmt.Sprintf("%s(%s/(%d), %s - (%s / %d) * %d )", v.TypeName, name, div, name, name, div, div)
}

func JsImportFormat(name string) string {
	return strings.ReplaceAll(name, "/", ".")
}

func JsImportName(name string) string {
	_, base := filepath.Split(name)
	return base
}

func JsInsFormat(data []NameAndType) string {
	res := []string{}

	for i, v := range data {
		typeName := jsTypeName(v.TypeName)
		valueName := v.Name
		if valueName == "_" {
			valueName = fmt.Sprintf("p%d", i+1)
		}
		res = append(res, fmt.Sprintf("%s: %s", valueName, typeName))
	}

	return strings.Join(res, ", ")
}

func JsOutsFormat(data []NameAndType) string {
	if data == nil {
		return ""
	}
	if len(data) == 1 {
		return ": " + jsTypeName(data[0].TypeName)
	}

	res := []string{}

	for _, v := range data {
		typeName := jsTypeName(v.TypeName)
		res = append(res, typeName)
	}

	if res[len(res)-1] == "error" {
		res = res[:len(res)-1]
	}

	if len(res) == 1 {
		return ": " + res[0]
	}

	return ": [" + strings.Join(res, ", ") + "]"
}

func GetRealElement(v gotype.Type) gotype.Type {
	kind := v.Kind()
	if kind == gotype.Ptr || kind == gotype.Slice || kind == gotype.Map {
		return GetRealElement(v.Elem())
	}
	return v
}

func GoFuncImportsPkgs(fn FuncDefine, pkgMap map[string]bool) {
	if len(fn.Outs) < 2 {
		return
	}

	for _, v := range fn.Ins {
		typeName := v.TypeName
		if strings.Contains(typeName, ".") {
			r := GetRealElement(v.Type)
			pkgName := r.PkgPath()
			pkgMap[pkgName] = true
		}
	}
}

func GoFuncFormat(fn FuncDefine) string {
	return "origin." + fn.Name

	if len(fn.Outs) < 2 {
		return "origin." + fn.Name
	}
	outs := []string{}
	for i := range fn.Outs {
		outs = append(outs, fmt.Sprintf("_fnoi_%d", i+1))
	}
	outstr := strings.Join(outs, ", ")

	ins := []string{}
	inValues := []string{}
	for _, v := range fn.Ins {
		typeName := v.TypeName
		if ast.IsExported(typeName) {
			typeName = "origin." + typeName
		}
		if v.Type.Kind() == gotype.Ptr {

		}
		ins = append(ins, fmt.Sprintf("%s %s", v.Name, typeName))
		inValues = append(inValues, v.Name)
	}
	instr := strings.Join(ins, ", ")
	invstr := strings.Join(inValues, ", ")

	return fmt.Sprintf(`func(%s) []any {
		%s := origin.%s(%s)
		return []any{%s}
	}`, instr, outstr, fn.Name, invstr, outstr)
}

func GoOutFormat(data []NameAndType) string {
	return ""
}

var (
	arrayRegxp = regexp.MustCompile(`\[.*\]`)
)

func jsTypeName(typeName string) string {
	if strings.Contains(typeName, "{") {
		return "object"
	}
	if strings.Contains(typeName, "func") {
		return "object"
	}
	if strings.Contains(typeName, "_") {
		return "any"
	}

	if strings.Contains(typeName, "chan") {
		typeName = strings.ReplaceAll(typeName, "<-", "")
		typeName = strings.ReplaceAll(typeName, "->", "")
		typeName = strings.ReplaceAll(typeName, "chan", "chanOf")
		typeName = strings.ReplaceAll(typeName, " ", "")
	}

	typeName = strings.ReplaceAll(typeName, "*", "")
	typeName = strings.ReplaceAll(typeName, "interface{}", "any")

	if strings.Contains(typeName, "[") {
		sliceExt := ""
		typeName = arrayRegxp.ReplaceAllStringFunc(typeName, func(s string) string {
			sliceExt += "[]"
			return ""
		})

		typeName = typeName + sliceExt
	}

	return typeName
}
