package web

import (
	"strconv"
)

type IntegerSchema struct {
	Type        string `json:"type,omitempty"`
	Title       string `json:"title,omitempty"`
	Default     int    `json:"default"`
	Description string `json:"description,omitempty"`
	Minimum     int    `json:"minimum"`
	Maximum     int    `json:"maximum,omitempty"`
	*AnnotationSchema
}

func NewIntegerSchema() *IntegerSchema {
	return &IntegerSchema{
		Type: "integer",
	}
}

func (i *IntegerSchema) SetTitle(title, name string) *IntegerSchema {
	if title != "" {
		i.Title = title
	} else {
		i.Title = titler(name)
	}
	return i
}
func (i *IntegerSchema) SetDefault(value string) *IntegerSchema {
	val, err := strconv.Atoi(value)
	if err == nil {
		i.Default = val
	}
	return i
}
func (i *IntegerSchema) SetDescription(value string) *IntegerSchema {
	i.Description = titler(value)
	return i
}

func (i *IntegerSchema) SetMinimum(value int) *IntegerSchema {
	i.Minimum = value
	return i
}
func (i *IntegerSchema) SetMaximum(value int) *IntegerSchema {
	if value > 0 {
		i.Maximum = value
	}
	return i
}
func (i *IntegerSchema) SetAnnotations(value *AnnotationSchema) *IntegerSchema {
	i.AnnotationSchema = value
	return i
}
