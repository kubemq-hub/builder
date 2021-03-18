package web

type MapSchema struct {
	Type  string `json:"type,omitempty"`
	Title string `json:"title"`
	Items struct {
		Type       string      `json:"type"`
		Required   []string    `json:"required"`
		Properties *OrderedMap `json:"properties"`
	} `json:"items"`

	*AnnotationSchema
}

func NewMapSchema() *MapSchema {
	m := &MapSchema{
		Type:  "array",
		Title: "",
		Items: struct {
			Type       string      `json:"type"`
			Required   []string    `json:"required"`
			Properties *OrderedMap `json:"properties"`
		}{
			Type:       "object",
			Required:   []string{"key", "value"},
			Properties: NewOrderedMap(),
		},
		AnnotationSchema: nil,
	}
	m.Items.Properties.Set("key", NewStringSchema().SetTitle("Set Key", ""))
	m.Items.Properties.Set("value", NewStringSchema().SetTitle("Set Value", ""))
	return m
}
func (m *MapSchema) SetTitle(title, name string) *MapSchema {
	if title != "" {
		m.Title = title
	} else {
		m.Title = titler(name)
	}
	return m
}
func (m *MapSchema) SetAnnotations(value *AnnotationSchema) *MapSchema {
	m.AnnotationSchema = value
	return m
}
