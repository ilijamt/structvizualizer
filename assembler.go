package structvizualizer

import (
	"go/ast"
	"text/template"

	"strings"

	"io"

	"fmt"
	"reflect"

	"github.com/davecgh/go-spew/spew"
)

type Assembler struct {
	GraphName   string
	Objects     map[string]*Object
	Files       map[string]*ast.File
	Connections map[string]ObjectConnection
}

func NewAssembler(name string) *Assembler {
	o := &Assembler{
		GraphName:   name,
		Objects:     make(map[string]*Object),
		Files:       make(map[string]*ast.File),
		Connections: make(map[string]ObjectConnection),
	}
	return o
}

func (o *Assembler) Parse() error {
	for filename, file := range o.Files {
		for _, decl := range file.Decls {
			o.parseDeclaration(filename, decl)
		}
	}
	o.buildConnections()
	return nil
}

func (o *Assembler) parseDeclaration(file string, decl ast.Decl) error {

	switch d := decl.(type) {
	case *ast.GenDecl:
		for _, spec := range d.Specs {
			o.parseSpec(file, spec)
		}
	default:
		//spew.Dump(d)
	}

	return nil
}

func (o *Assembler) parseSpec(file string, spec ast.Spec) error {
	switch s := spec.(type) {
	case *ast.TypeSpec:
		name := s.Name.String()
		o.Objects[name] = NewObject(name, file)

		switch sp := s.Type.(type) {
		case *ast.StructType:
			for _, field := range sp.Fields.List {
				o.parseStructField(name, field)
			}

		case *ast.Ident:
			field := NewObjectField()
			field.Type = new(string)
			field.Primitive = true
			*field.Type = sp.String()
			o.Objects[name].AddLabel(field)

		case *ast.InterfaceType:
			for _, field := range sp.Methods.List {
				o.parseStructField(name, field)
			}

		default:
			fmt.Printf("(s.Type.(type)) %s Ignoring %s type\n", s.Name.String(), reflect.TypeOf(sp))
		}

	case *ast.ImportSpec:
		// ignore import

	case *ast.ValueSpec:
		//fmt.Printf("(*ast.ValueSpec) %#v \n", s)

	default:
		fmt.Printf("(s := spec.(type)) Ignoring %s type\n", reflect.TypeOf(s))
	}

	return nil
}

func (o *Assembler) parseStructField(structName string, field *ast.Field) error {

	obj := NewObjectField()

	if field.Tag != nil {
		obj.Tag = new(string)
		*obj.Tag = field.Tag.Value
	}

	if len(field.Names) > 0 {
		obj.Name = new(string)
		*obj.Name = field.Names[0].Name
	}

	obj.Type = new(string)

	switch t := field.Type.(type) {
	case *ast.SelectorExpr:
		*obj.Type = t.X.(*ast.Ident).Name + t.Sel.Name
	case *ast.Ident:
		*obj.Type = t.String()
	case *ast.InterfaceType:
		*obj.Type = "interface{}"
	case *ast.ArrayType:
		obj.IsArray = true
		switch e := t.Elt.(type) {
		case *ast.Ident:
			*obj.Type = e.Name
		case *ast.SelectorExpr:
			*obj.Type = e.X.(*ast.Ident).Name + e.Sel.Name
		default:
		}
	case *ast.StarExpr:

		switch x := t.X.(type) {
		case *ast.SelectorExpr:
			*obj.Type = x.X.(*ast.Ident).Name + x.Sel.Name
		case *ast.Ident:
			*obj.Type = x.Name
		default:
			fmt.Printf("(*ast.StarExpr) Ignoring %s type\n", reflect.TypeOf(x))
		}
	case *ast.MapType:
		var keyType string
		var valType string

		switch e := t.Key.(type) {
		case *ast.Ident:
			keyType = e.Name
		default:
			keyType = "Unknown"
		}

		switch e := t.Value.(type) {
		case *ast.Ident:
			valType = e.Name
		default:
			valType = "Unknown"
		}
		*obj.Type = fmt.Sprintf("map[%s]%s", keyType, valType)

	case *ast.FuncType:
		obj.Function = true
		*obj.Type = "func"
	default:
		fmt.Printf("(field.Type.(type)) Ignoring %s type\n", reflect.TypeOf(t))
	}

	o.Objects[structName].AddField(obj)

	return nil
}

func (o *Assembler) Add(file string, astFile *ast.File) error {
	o.Files[file] = astFile
	return nil
}

func (o *Assembler) Dump() {
	spew.Dump(o.Objects)
}

func (o *Assembler) buildConnections() {

	for _, v := range o.Objects {

		// loop over all the fields to build the connection map
		for _, field := range v.Fields {
			_, referenced := o.Objects[*field.Type]

			embedded := field.IsEmbedded()

			if embedded {
				fieldName := field.GetName()
				label := "embedded"

				obj := NewObjectConnection(v.Name, fieldName, label)
				o.Connections[obj.Hash()] = obj
			} else if field.IsPrimitive() {
				fieldName := field.GetName()
				label := "type"
				obj := NewObjectConnection(v.Name, fieldName, label)
				o.Connections[obj.Hash()] = obj
			} else if referenced {
				label := ""
				if field.IsArray {
					label = "[]" + label
				}
				label += field.GetName()

				obj := NewObjectConnection(v.Name, *field.Type, label)
				o.Connections[obj.Hash()] = obj

			}

		}

	}

}

func (o *Assembler) Template(output io.Writer) {

	tmpl := template.New(o.GraphName).Funcs(template.FuncMap{"join": strings.Join})

	tmpl, err := tmpl.Parse(Template)
	if err != nil {
		panic(err)
	}
	tmpl.Execute(output, o)
}
