package main

import (
	"bytes"
	"embed"
	"fmt"
	"go/ast"
	"go/format"
	"io/ioutil"
	"path/filepath"
	"strings"
	"text/template"

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
	Name     string
	Doc      string
	Comment  string
	TypeName string
	Type     gotype.Type
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
	imp := gotype.NewImporter()
	pkgInfos, err := imp.Import(pkgPath, "")
	if err != nil {
		return err
	}

	values := []ValueDefine{}
	types := []TypeDefine{}
	funcs := []FuncDefine{}

	num := pkgInfos.NumChild()

	for i := 0; i < num; i++ {
		v := pkgInfos.Child(i)
		if ast.IsExported(v.Name()) {
			switch v.Kind() {
			case gotype.Declaration:
				decl := v.Declaration()
				switch decl.Kind() {
				case gotype.Func:
					ImportFunc(v, &funcs)
				default:
					ImportValues(v, &values)
				}
			case gotype.Scope, gotype.Invalid:
			default:
				ImportTypes(v, &types)
			}
		}
	}

	buf := bytes.NewBuffer(nil)

	pkgPath = strings.Trim(pkgPath, "/")

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
	}

	templateFile.ExecuteTemplate(buf, "tempGo", tmplData)

	code := buf.String()

	codeFormated, err := format.Source([]byte(code))

	if err != nil {
		ioutil.WriteFile(dstGoFile, []byte(code), 0755)
		fmt.Println("Err", err)
		return err
	}

	ioutil.WriteFile(dstGoFile, codeFormated, 0755)

	bufJs := bytes.NewBuffer(nil)
	jsTemplate.ExecuteTemplate(bufJs, "tempJS", tmplData)
	ioutil.WriteFile(dstJSFile, bufJs.Bytes(), 0755)

	return nil
}

func ImportFunc(v gotype.Type, fnList *[]FuncDefine) {
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

	(*fnList) = append(*fnList, fn)
}

func ImportTypes(v gotype.Type, typeList *[]TypeDefine) {
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
					ImportValues(para, &fields)
				}
			}
			tp.Fields = fields
		}
	}

	if v.NumMethod() > 0 {
		methods := []FuncDefine{}
		for i := 0; i < v.NumMethod(); i++ {
			fn := v.Method(i)
			if fn.Kind() == gotype.Declaration {
				if ast.IsExported(fn.Name()) {
					ImportFunc(fn, &methods)
				}
			} else {
				fmt.Println(fn.Kind())
			}

		}
		tp.Funcs = methods
	}

	(*typeList) = append((*typeList), tp)
}

func ImportValues(v gotype.Type, valueList *[]ValueDefine) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("PANIC:", v.Kind(), "|", v.Name(), "|", v)
		}
	}()

	typeName := ""

	if v.Kind() == gotype.Declaration {
		typeName = v.Declaration().String()
	} else {
		name := v.Name()
		s := v.String()
		typeName = strings.TrimSpace(strings.Replace(s, name, "", 1))
	}

	tp := ValueDefine{
		Name:     v.Name(),
		Doc:      comment(v.Doc().Text()),
		Comment:  comment(v.Comment().Text()),
		Type:     v,
		TypeName: typeName,
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
	templateFile.Parse(goTemp)

	jsTemplate.Funcs(template.FuncMap{
		"jsInsFormat":  JsInsFormat,
		"jsOutsFormat": JsOutsFormat,
		"jsTypeName":   jsTypeName,
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

func JsInsFormat(data []NameAndType) string {
	res := []string{}

	for _, v := range data {
		typeName := jsTypeName(v.TypeName)
		res = append(res, fmt.Sprintf("%s: %s", v.Name, typeName))
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
	return ": [" + strings.Join(res, ", ") + "]"
}

func jsTypeName(typeName string) string {
	typeName = strings.ReplaceAll(typeName, "*", "")
	typeName = strings.ReplaceAll(typeName, "interface{}", "any")
	if strings.Contains(typeName, "[]") {
		typeName = strings.ReplaceAll(typeName, "[]", "")
		typeName = typeName + "[]"
	}
	return typeName
}
