package web

import (
	"github.com/kubemq-hub/builder/connector/common"
)

func ConvertToJsonSchemaList(connectors common.Connectors) (*JsonSchemaList, error) {
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
