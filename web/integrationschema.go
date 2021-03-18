package web

import (
	"encoding/json"
	"fmt"
	"github.com/kubemq-hub/builder/connector/common"
	"sort"
)

type IntegrationsSchema struct {
	Integrations  []*KindMetadata `json:"integrations"`
	KubemqSources *JsonSchemaList `json:"kubemq_sources"`
	KubemqTargets *JsonSchemaList `json:"kubemq_targets"`
}

func NewIntegrationSchema() *IntegrationsSchema {
	return &IntegrationsSchema{}
}
func propertyToJsonComponent(p *common.Property) (string, interface{}, error) {
	switch p.Kind {
	case "string":
		return p.Name, NewStringSchema().
				SetTitle(p.Title, p.Name).
				SetDescription(p.Description).
				SetEnum(p.Options...).
				SetDefault(p.Default),
			nil
	case "condition":
		conditional := NewSelectConditionSchema().SetTitle(p.Title, p.Name)
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
				SetTitle(p.Title, p.Name).
				SetDescription(p.Description).
				SetDefault(p.Default).
				SetMinimum(p.Min).
				SetMaximum(p.Max),
			nil
	case "multilines":
		return p.Name, NewStringSchema().
				SetTitle(p.Title, p.Name).
				SetDescription(p.Description).
				SetDefault(p.Default).
				SetAnnotations(NewAnnotationSchema().SetDisplay("textarea")),
			nil
	case "bool":
		return p.Name, NewBooleanSchema().
				SetTitle(p.Title, p.Name).
				SetDescription(p.Description).
				SetDefault(p.Default),
			nil
	case "map":
		m := NewMapSchema()
		title := titler(p.Name)
		if p.Title != "" {
			title = p.Title
		}
		m.SetTitle(fmt.Sprintf("Add %s Key Value Pairs", title), "")
		return p.Name, m, nil
	case "null":
		return "", nil, nil
	}
	return "", nil, fmt.Errorf("kind not found: %s", p.Kind)
}
func (i *IntegrationsSchema) toJsonSchemaList(connectors common.Connectors) (*JsonSchemaList, error) {
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
func (i *IntegrationsSchema) toKindMetadata(group string, connectors common.Connectors) ([]*KindMetadata, error) {
	var list []*KindMetadata
	for _, connector := range connectors {
		j := NewJsonSchemaItem().SetKind(connector.Kind)
		j.AnnotationSchema = NewAnnotationSchema().SetClass("vjsf")
		for _, property := range connector.Properties {
			key, val, err := propertyToJsonComponent(property)
			if err != nil {
				return nil, err
			}
			j.SetItemProperties(key, val, property.Must)
		}
		list = append(list, NewKindMetadata().
			SetType(group).
			SetKind(connector.Kind).
			SetName(connector.Name).
			SetCategory(connector.Category).
			SetProvider(connector.Provider).
			SetTags(connector.Tags).
			SetSchema(j))
	}
	return list, nil
}
func (i *IntegrationsSchema) Load(sources, targets string) (*IntegrationsSchema, error) {
	targetManifest, err := common.LoadManifestFromFile(targets)
	if err != nil {
		return nil, err
	}

	targetsList, err := i.toKindMetadata("targets", targetManifest.Targets)
	if err != nil {
		return nil, err
	}
	i.Integrations = append(i.Integrations, targetsList...)

	sourcesManifest, err := common.LoadManifestFromFile(sources)
	if err != nil {
		return nil, err
	}

	sourcesList, err := i.toKindMetadata("sources", sourcesManifest.Sources)
	if err != nil {
		return nil, err
	}
	i.Integrations = append(i.Integrations, sourcesList...)

	i.KubemqTargets, err = i.toJsonSchemaList(sourcesManifest.Targets)
	if err != nil {
		return nil, err
	}

	i.KubemqSources, err = i.toJsonSchemaList(targetManifest.Sources)
	if err != nil {
		return nil, err
	}
	sort.Slice(i.Integrations, func(t, j int) bool {
		return i.Integrations[t].Name < i.Integrations[j].Name

	})
	return i, nil
}

func (i *IntegrationsSchema) Marshal() []byte {
	data, _ := json.MarshalIndent(i, "", "\t")
	return data
}

func (i *IntegrationsSchema) String() string {
	return string(i.Marshal())
}
