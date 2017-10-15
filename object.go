package structvizualizer

import "fmt"

type Object struct {
	Name   string
	File   string
	Label  []string
	Fields map[string]ObjectField
}

func NewObject(name string, file string) *Object {
	return &Object{
		Name:   name,
		File:   file,
		Fields: make(map[string]ObjectField),
	}
}

func (o *Object) AddLabel(field ObjectField) {
	tp := "prop: "
	if field.IsEmbedded() {
		tp = "embd: "
	}
	if field.IsPrimitive() {
		tp = ""
	}
	if field.IsFunction() {
		tp = "func: "
	}
	o.Label = append(o.Label, fmt.Sprintf("{%s%s}", tp, field.GetName()))
}

func (o *Object) AddField(field ObjectField) {
	o.Fields[field.GetName()] = field
	o.AddLabel(field)
}
