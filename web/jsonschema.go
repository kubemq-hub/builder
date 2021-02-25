package web

type JsonSchemaItem struct {
	*BaseSchema
	If             *BaseSchema `json:"if,omitempty"`
	Then           *BaseSchema `json:"then,omitempty"`
	conditionSetup bool
}

func NewJsonSchemaItem() *JsonSchemaItem {
	return &JsonSchemaItem{
		BaseSchema: NewBaseSchema(),
		If:         nil,
		Then:       nil,
	}
}
func (j *JsonSchemaItem) setupCondition() {
	j.Properties.Set("setDefaults",
		NewBooleanSchema().
			SetTitle("Set Defaults Properties").
			SetDefault("true").
			SetAnnotations(NewAnnotationSchema().
				SetDisplay("checkbox")))
	j.If = NewBaseSchema().SetRequired("setDefaults")
	j.If.Properties.Set("setDefaults", struct {
		Const bool `json:"const"`
	}{
		Const: false,
	})

	j.Then = NewBaseSchema()
	j.conditionSetup = true
}
func (j *JsonSchemaItem) SetItemProperties(key string, value interface{}, isRequired bool) *JsonSchemaItem {
	if isRequired {
		j.SetRequired(key)
		if j.Properties == nil {
			j.Properties = NewOrderedMap()
		}
		j.Properties.Set(key, value)
	} else {
		if !j.conditionSetup {
			j.setupCondition()

		}
		j.Then.Properties.Set(key, value)
	}

	return j
}

type JsonSchemaList struct {
	Title string            `json:"title"`
	Type  string            `json:"type"`
	OneOf []*JsonSchemaItem `json:"oneOf"`
}

func NewJsonSchemaList() *JsonSchemaList {
	return &JsonSchemaList{
		Title: "schema",
		Type:  "object",
		OneOf: nil,
	}
}
func (j *JsonSchemaList) AddItem(value *JsonSchemaItem) *JsonSchemaList {
	j.OneOf = append(j.OneOf, value)
	return j
}
