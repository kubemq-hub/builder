package common

import (
	"fmt"
	"github.com/kubemq-hub/builder/pkg/utils"
	"github.com/kubemq-hub/builder/survey"
)

type Binding struct {
	Name              string            `json:"name" yaml:"name"`
	Source            Spec              `json:"source" yaml:"source"`
	Target            Spec              `json:"target" yaml:"target"`
	Properties        map[string]string `json:"properties" yaml:"properties"`
	SourceSpec        string            `json:"-" yaml:"-"`
	TargetSpec        string            `json:"-" yaml:"-"`
	PropertiesSpec    string            `json:"-" yaml:"-"`
	loadedOptions     DefaultOptions
	targetsList       Connectors
	sourcesList       Connectors
	takenBindingNames []string
	defaultName       string
	isEditMode        bool
}

func NewBinding(defaultName string) *Binding {
	return &Binding{
		Name:              "",
		Source:            Spec{},
		Target:            Spec{},
		Properties:        map[string]string{},
		SourceSpec:        "",
		TargetSpec:        "",
		PropertiesSpec:    "",
		loadedOptions:     nil,
		targetsList:       nil,
		sourcesList:       nil,
		takenBindingNames: nil,
		defaultName:       defaultName,
		isEditMode:        false,
	}
}
func (b *Binding) SetDefaultOptions(value DefaultOptions) *Binding {
	b.loadedOptions = value
	return b
}
func (b *Binding) Clone() *Binding {
	newBinding := &Binding{
		Name:              b.Name,
		Source:            b.Source.Clone(),
		Target:            b.Target.Clone(),
		Properties:        map[string]string{},
		SourceSpec:        "",
		TargetSpec:        "",
		PropertiesSpec:    "",
		loadedOptions:     nil,
		targetsList:       nil,
		sourcesList:       nil,
		takenBindingNames: nil,
		defaultName:       "",
		isEditMode:        false,
	}
	for key, val := range b.Properties {
		newBinding.Properties[key] = val
	}
	return newBinding
}
func (b *Binding) SetTargetsList(value Connectors) *Binding {
	b.targetsList = value
	return b
}
func (b *Binding) SetSourcesList(value Connectors) *Binding {
	b.sourcesList = value
	return b
}
func (b *Binding) SetEditMode(value bool) *Binding {
	b.isEditMode = value
	return b
}
func (b *Binding) SetTakenBindingNames(value []string) *Binding {
	b.takenBindingNames = value
	return b
}
func (b *Binding) SourceName() string {
	return b.Source.Name
}
func (b *Binding) TargetName() string {
	return b.Target.Name
}
func (b *Binding) askKind(kinds []string, currentKind string) (string, error) {
	defaultKind := ""
	if b.isEditMode {
		defaultKind = currentKind
	} else {
		defaultKind = kinds[0]
	}
	if defaultKind == "" {
		defaultKind = kinds[0]
	}
	val := ""
	err := survey.NewString().
		SetKind("string").
		SetName("kind").
		SetMessage("Select Connector Kind").
		SetDefault(defaultKind).
		SetOptions(kinds).
		SetHelp("Select Connector Kind").
		SetRequired(true).
		Render(&val)
	if err != nil {
		return "", err
	}
	return val, nil
}

func (b *Binding) confirmSource() bool {
	utils.Println(fmt.Sprintf(promptSourceConfirm, b.Source.String(sourceSpecTemplate)))
	val := true
	err := survey.NewBool().
		SetKind("bool").
		SetName("confirm-connection").
		SetMessage("Would you like save this configuration").
		SetDefault("true").
		SetRequired(true).
		Render(&val)
	if err != nil {
		return false
	}

	return val
}
func (b *Binding) confirmTarget() bool {
	utils.Println(fmt.Sprintf(promptTargetConfirm, b.Target.String(targetSpecTemplate)))
	val := true
	err := survey.NewBool().
		SetKind("bool").
		SetName("confirm-connection").
		SetMessage("Would you like save this configuration").
		SetDefault("true").
		SetRequired(true).
		Render(&val)
	if err != nil {
		return false
	}

	return val
}
func (b *Binding) confirmProperties(p *Properties) bool {
	utils.Println(fmt.Sprintf(promptPropertiesConfirm, p.String()))
	val := true
	err := survey.NewBool().
		SetKind("bool").
		SetName("confirm-connection").
		SetMessage("Would you like save this configuration").
		SetDefault("true").
		SetRequired(true).
		Render(&val)
	if err != nil {
		return false
	}
	if !val {
		utils.Println(promptPropertiesReconfigure)
	}
	return val
}
func (b *Binding) askSource(defaultName string) error {
	var err error
	sourceDefaultName := ""
	if b.isEditMode {
		sourceDefaultName = b.Source.Name
	} else {
		sourceDefaultName = defaultName
	}
	if b.Source.Name, err = NewName(sourceDefaultName).
		RenderSource(); err != nil {
		return err
	}
	var kinds []string
	sources := make(map[string]*Connector)
	for _, c := range b.sourcesList {
		kinds = append(kinds, c.Kind)
		sources[c.Kind] = c
	}
	if len(kinds) == 0 {
		return fmt.Errorf("no source connectors available")
	}

	if b.Source.Kind, err = b.askKind(kinds, b.Source.Kind); err != nil {
		return err
	}
	connector := sources[b.Source.Kind]
	if b.Source.Properties, err = connector.Render(b.loadedOptions); err != nil {
		return err
	}
	return nil
}
func (b *Binding) askTarget(defaultName string) error {
	var err error
	targetDefaultName := ""
	if b.isEditMode {
		targetDefaultName = b.Target.Name
	} else {
		targetDefaultName = defaultName
	}
	if b.Target.Name, err = NewName(targetDefaultName).
		RenderTarget(); err != nil {
		return err
	}
	var kinds []string
	targets := make(map[string]*Connector)
	for _, c := range b.targetsList {
		kinds = append(kinds, c.Kind)
		targets[c.Kind] = c
	}
	if len(kinds) == 0 {
		return fmt.Errorf("no targets connectors available")
	}

	if b.Target.Kind, err = b.askKind(kinds, b.Target.Kind); err != nil {
		return err
	}
	connector := targets[b.Target.Kind]
	if b.Target.Properties, err = connector.Render(b.loadedOptions); err != nil {
		return err
	}
	return nil
}

func (b *Binding) Render() (*Binding, error) {
	var err error
	defaultName := ""
	if b.isEditMode {
		defaultName = b.Name
	} else {
		defaultName = b.defaultName
	}
	if b.Name, err = NewName(defaultName).
		SetTakenNames(b.takenBindingNames).
		RenderBinding(); err != nil {
		return nil, err
	}
	utils.Println(promptSourceStart)
	for {
		if err := b.askSource(fmt.Sprintf("%s-source", b.Name)); err != nil {
			return nil, err
		}
		ok := b.confirmSource()
		if ok {
			b.SourceSpec = b.Source.String(sourceSpecTemplate)
			break
		}
	}
	utils.Println(promptTargetStart)
	for {
		if err := b.askTarget(fmt.Sprintf("%s-target", b.Name)); err != nil {
			return nil, err
		}
		ok := b.confirmTarget()
		if ok {
			b.TargetSpec = b.Target.String(targetSpecTemplate)
			break
		}
	}
	utils.Println(promptBindingComplete)
	for {
		p := NewProperties()
		if b.Properties, err = p.
			Render(); err != nil {
			return nil, err
		}
		if len(b.Properties) == 0 {
			break
		}
		ok := b.confirmProperties(p)
		if ok {
			b.PropertiesSpec = p.String()
			break
		}
	}

	return b, nil
}

func (b *Binding) String() string {
	tpl := utils.NewTemplate(bindingTemplate, b)
	bnd, err := tpl.Get()
	if err != nil {
		return fmt.Sprintf("error rendring binding spec,%s", err.Error())
	}
	return string(bnd)
}
