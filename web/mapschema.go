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
	m.Items.Properties.Set("key", NewStringSchema().SetTitle("Set Key"))
	m.Items.Properties.Set("value", NewStringSchema().SetTitle("Set Value"))
	return m
}
func (s *MapSchema) SetTitle(value string) *MapSchema {
	s.Title = value
	return s
}
func (s *MapSchema) SetAnnotations(value *AnnotationSchema) *MapSchema {
	s.AnnotationSchema = value
	return s
}
