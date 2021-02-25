package web

type BaseSchema struct {
	Type        string      `json:"type,omitempty"`
	Title       string      `json:"title,omitempty"`
	Description string      `json:"description,omitempty"`
	Required    []string    `json:"required,omitempty"`
	Properties  *OrderedMap `json:"properties,omitempty"`
	*AnnotationSchema
}

func NewBaseSchema() *BaseSchema {
	return &BaseSchema{
		Properties: NewOrderedMap(),
	}
}

func (b *BaseSchema) SetType(value string) *BaseSchema {
	b.Type = value
	return b
}

func (b *BaseSchema) SetTitle(value string) *BaseSchema {
	b.Title = value
	return b
}

func (b *BaseSchema) SetDescription(value string) *BaseSchema {
	b.Description = value
	return b
}

func (b *BaseSchema) SetRequired(values ...string) *BaseSchema {
	b.Required = append(b.Required, values...)
	return b
}

func (b *BaseSchema) SetProperties(key string, value interface{}) *BaseSchema {
	if b.Properties == nil {
		b.Properties = NewOrderedMap()
	}
	b.Properties.Set(key, value)
	return b
}
func (b *BaseSchema) SetAnnotations(value *AnnotationSchema) *BaseSchema {
	b.AnnotationSchema = value
	return b
}
