package web

import (
	"strconv"
)

type BooleanSchema struct {
	Type        string `json:"type,omitempty"`
	Title       string `json:"title,omitempty"`
	Default     bool   `json:"default"`
	Description string `json:"description,omitempty"`
	*AnnotationSchema
}

func NewBooleanSchema() *BooleanSchema {
	return &BooleanSchema{
		Type: "boolean",
	}
}

func (b *BooleanSchema) SetTitle(value string) *BooleanSchema {
	b.Title = titler(value)
	return b
}
func (b *BooleanSchema) SetDefault(value string) *BooleanSchema {
	b.Default, _ = strconv.ParseBool(value)
	return b
}
func (b *BooleanSchema) SetDescription(value string) *BooleanSchema {
	b.Description = titler(value)
	return b
}

func (b *BooleanSchema) SetAnnotations(value *AnnotationSchema) *BooleanSchema {
	b.AnnotationSchema = value
	return b
}
