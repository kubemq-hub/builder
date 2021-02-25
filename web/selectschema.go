package web

type SubSelect struct {
	Title      string      `json:"title,omitempty"`
	Properties *OrderedMap `json:"properties"`
}

func NewSubSelect(name string) *SubSelect {
	ss := &SubSelect{
		Title:      titler(name),
		Properties: NewOrderedMap(),
	}
	ss.Properties.Set("key", struct {
		Type  string `json:"type"`
		Const string `json:"const"`
	}{
		Type:  "string",
		Const: name,
	})
	return ss
}

type SelectConditionSchema struct {
	Type  string       `json:"type,omitempty"`
	Title string       `json:"title,omitempty"`
	OneOf []*SubSelect `json:"oneOf,omitempty"`
	*AnnotationSchema
}

func NewSelectConditionSchema() *SelectConditionSchema {
	return &SelectConditionSchema{
		Type: "object",
	}
}
func (s *SelectConditionSchema) SetTitle(value string) *SelectConditionSchema {
	s.Title = titler(value)
	return s
}
func (s *SelectConditionSchema) SetOneOf(value string, properties *OrderedMap) *SelectConditionSchema {
	ss := NewSubSelect(value)
	keys := properties.keys
	for _, k := range keys {
		v, _ := properties.Get(k)
		ss.Properties.Set(k, v)
	}
	s.OneOf = append(s.OneOf, ss)
	return s
}

func (s *SelectConditionSchema) SetAnnotations(value *AnnotationSchema) *SelectConditionSchema {
	s.AnnotationSchema = value
	return s
}
