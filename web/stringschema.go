package web

type StringSchema struct {
	Type        string   `json:"type,omitempty"`
	Title       string   `json:"title,omitempty"`
	Default     string   `json:"default,omitempty"`
	Description string   `json:"description,omitempty"`
	Enum        []string `json:"enum,omitempty"`
	*AnnotationSchema
}

func NewStringSchema() *StringSchema {
	return &StringSchema{
		Type: "string",
	}
}

func (s *StringSchema) SetTitle(value string) *StringSchema {
	s.Title = titler(value)
	return s
}
func (s *StringSchema) SetDefault(value string) *StringSchema {
	s.Default = value
	return s
}
func (s *StringSchema) SetDescription(value string) *StringSchema {
	s.Description = titler(value)
	return s
}
func (s *StringSchema) SetEnum(values ...string) *StringSchema {
	s.Enum = append(s.Enum, values...)
	return s
}
func (s *StringSchema) SetAnnotations(value *AnnotationSchema) *StringSchema {
	s.AnnotationSchema = value
	return s
}
