package common

import (
	"fmt"
	"github.com/kubemq-hub/builder/survey"
)

type Connector struct {
	Kind          string      `json:"kind"`
	Description   string      `json:"description"`
	Properties    []*Property `json:"properties"`
	loadedOptions DefaultOptions
	values        map[string]string
}

func NewConnector() *Connector {
	return &Connector{
		Kind:          "",
		Description:   "",
		Properties:    nil,
		loadedOptions: nil,
	}
}

func (c *Connector) SetKind(value string) *Connector {
	c.Kind = value
	return c
}

func (c *Connector) SetDescription(value string) *Connector {
	c.Description = value
	return c
}
func (c *Connector) AddProperty(value *Property) *Connector {
	c.Properties = append(c.Properties, value)
	return c
}

func (c *Connector) askString(p *Property) error {
	val := ""
	options := c.loadedOptions[fmt.Sprintf("%s/%s", c.Kind, p.Name)]
	err := survey.NewString().
		SetKind("string").
		SetName(p.Name).
		SetMessage(p.Description).
		SetDefault(p.Default).
		SetOptions(options).
		SetHelp(p.Description).
		SetRequired(p.Must).
		Render(&val)
	if err != nil {
		return err
	}
	if val != "" {
		c.values[p.Name] = val
	}
	return nil
}
func (c *Connector) askInt(p *Property) error {
	val := 0
	err := survey.NewInt().
		SetKind("int").
		SetName(p.Name).
		SetMessage(p.Description).
		SetDefault(p.Default).
		SetHelp(p.Description).
		SetRequired(p.Must).
		SetRange(p.Min, p.Max).
		Render(&val)
	if err != nil {
		return err
	}
	c.values[p.Name] = fmt.Sprintf("%d", val)
	return nil
}
func (c *Connector) askBool(p *Property) error {
	val := false
	err := survey.NewBool().
		SetKind("bool").
		SetName(p.Name).
		SetMessage(p.Description).
		SetDefault(p.Default).
		Render(&val)
	if err != nil {
		return err
	}
	c.values[p.Name] = fmt.Sprintf("%t", val)
	return nil
}

func (c *Connector) Render(options DefaultOptions) (map[string]string, error) {
	c.values = map[string]string{}
	c.loadedOptions = options
	for _, p := range c.Properties {
		switch p.Kind {
		case "string":
			if err := c.askString(p); err != nil {
				return nil, err
			}
		case "int":
			if err := c.askInt(p); err != nil {
				return nil, err
			}
		case "bool":
			if err := c.askBool(p); err != nil {
				return nil, err
			}
		}
	}
	return c.values, nil
}
