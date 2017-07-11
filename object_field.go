package structvizualizer

type ObjectField struct {
	Name    *string
	Type    *string
	Tag     *string
	IsArray bool
}

func (o ObjectField) GetName() string {
	if o.IsEmbedded() {
		return *o.Type
	}
	return *o.Name
}

func (o ObjectField) IsEmbedded() bool {
	return o.Name == nil && o.Type != nil
}

func NewObjectField() ObjectField {
	return ObjectField{}
}
