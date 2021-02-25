package web

type Component struct {
	Type     string          `json:"type"`
	Name     string          `json:"name"`
	Category string          `json:"category"`
	Provider string          `json:"provider"`
	Schema   *JsonSchemaItem `json:"schema"`
}

func NewComponent() *Component {
	return &Component{}
}
func (c *Component) SetType(value string) *Component {
	c.Type = value
	return c
}
func (c *Component) SetName(value string) *Component {
	c.Name = value
	return c
}
func (c *Component) SetCategory(value string) *Component {
	c.Category = value
	return c
}
func (c *Component) SetProvider(value string) *Component {
	c.Provider = value
	return c
}
func (c *Component) SetSchema(value *JsonSchemaItem) *Component {
	c.Schema = value
	return c
}
