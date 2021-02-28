package web

import "strings"

type KindMetadata struct {
	Kind     string          `json:"kind"`
	Type     string          `json:"type"`
	Name     string          `json:"name"`
	Category string          `json:"category"`
	Provider string          `json:"provider"`
	Schema   *JsonSchemaItem `json:"schema"`
}

func NewKindMetadata() *KindMetadata {
	return &KindMetadata{}
}
func (s *KindMetadata) SetType(value string) *KindMetadata {
	s.Type = value
	return s
}
func (s *KindMetadata) SetKind(value string) *KindMetadata {
	s.Kind = value
	parts := strings.Split(value, ".")
	if len(parts) == 1 {
		s.Name = titler(parts[0])
		s.Category = "General"
		s.Provider = "Local"
		return s
	}
	if len(parts) == 2 {
		s.Name = titler(parts[1])
		s.Category = titler(parts[0])
		s.Provider = "Local"
		return s
	}
	if len(parts) == 3 {
		s.Name = titler(parts[2])
		s.Category = titler(parts[1])
		s.Provider = titler(parts[0])
		return s
	}

	return s
}
func (s *KindMetadata) SetSchema(value *JsonSchemaItem) *KindMetadata {
	s.Schema = value
	return s
}
