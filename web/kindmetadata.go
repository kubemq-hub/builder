package web

type KindMetadata struct {
	Kind     string          `json:"kind"`
	Type     string          `json:"type"`
	Name     string          `json:"name"`
	Category string          `json:"category"`
	Provider string          `json:"provider"`
	Tags     []string        `json:"tags"`
	Schema   *JsonSchemaItem `json:"schema"`
}

func NewKindMetadata() *KindMetadata {
	return &KindMetadata{}
}
func (s *KindMetadata) SetType(value string) *KindMetadata {
	s.Type = value
	return s
}

func (s *KindMetadata) SetName(value string) *KindMetadata {
	s.Name = value
	return s
}
func (s *KindMetadata) SetProvider(value string) *KindMetadata {
	s.Provider = value
	return s
}
func (s *KindMetadata) SetCategory(value string) *KindMetadata {
	s.Category = value
	return s
}
func (s *KindMetadata) SetTags(value []string) *KindMetadata {
	s.Tags = value
	return s
}

func (s *KindMetadata) SetKind(value string) *KindMetadata {
	s.Kind = value
	return s

}
func (s *KindMetadata) SetSchema(value *JsonSchemaItem) *KindMetadata {
	s.Schema = value
	return s
}
