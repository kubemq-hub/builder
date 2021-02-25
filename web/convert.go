package web

import (
	"fmt"
	"github.com/kubemq-hub/builder/connector/common"
)

func ConvertToJsonSchema(connectors common.Connectors) (*JsonSchemaList, error) {
	schema := NewJsonSchemaList()
	for _, connector := range connectors {
		j := NewJsonSchemaItem()
		j.BaseSchema.SetType("object").SetTitle(titler(connector.Kind))
		j.BaseSchema.Properties.Set("kind", struct {
			Type  string `json:"type"`
			Const string `json:"const"`
		}{
			Type:  "string",
			Const: connector.Kind,
		})
		for _, property := range connector.Properties {
			key, val, err := propertyToJsonComponent(property)
			if err != nil {
				return nil, err
			}
			j.SetItemProperties(key, val, property.Must)
		}
		schema.AddItem(j)
	}
	return schema, nil
}

func propertyToJsonComponent(p *common.Property) (string, interface{}, error) {
	switch p.Kind {
	case "string":
		return p.Name, NewStringSchema().
				SetTitle(p.Name).
				SetDescription(p.Description).
				SetEnum(p.Options...).
				SetDefault(p.Default),
			nil
	case "condition":
		conditional := NewSelectConditionSchema().SetTitle(p.Name)
		for name, properties := range p.Conditional {
			subProperties := NewOrderedMap()
			for _, property := range properties {
				key, val, err := propertyToJsonComponent(property)
				if err != nil {
					return "", nil, err
				}
				subProperties.Set(key, val)
			}
			conditional.SetOneOf(name, subProperties)
		}
		return p.Name, conditional, nil
	case "int":
		return p.Name, NewIntegerSchema().
				SetTitle(p.Name).
				SetDescription(p.Description).
				SetDefault(p.Default).
				SetMinimum(p.Min).
				SetMaximum(p.Max),
			nil
	case "multilines":
		return p.Name, NewStringSchema().
				SetTitle(p.Name).
				SetDescription(p.Description).
				SetDefault(p.Default).
				SetAnnotations(NewAnnotationSchema().SetDisplay("textarea")),
			nil
	case "bool":
		return p.Name, NewBooleanSchema().
				SetTitle(p.Name).
				SetDescription(p.Description).
				SetDefault(p.Default),
			nil
	case "map":
		m := NewMapSchema()
		m.SetTitle(fmt.Sprintf("Add %s Key Value Pairs", titler(p.Name)))
		return p.Name, m, nil
	case "null":
		return "", nil, nil
	}
	return "", nil, fmt.Errorf("kind not found: %s", p.Kind)
}
