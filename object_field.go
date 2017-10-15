package structvizualizer

type ObjectField struct {
	Name      *string
	Type      *string
	Tag       *string
	IsArray   bool
	Primitive bool
	Function  bool
}

func (o ObjectField) GetName() string {
	if o.IsEmbedded() || o.IsPrimitive() {
		return *o.Type
	}
	return *o.Name
}

func (o ObjectField) IsPrimitive() bool {
	return o.Primitive
}

func (o ObjectField) IsFunction() bool {
	return o.Function
}

func (o ObjectField) IsEmbedded() bool {
	return o.Name == nil && o.Type != nil && !o.IsPrimitive() && !o.IsFunction()
}

func NewObjectField() ObjectField {
	return ObjectField{}
}
