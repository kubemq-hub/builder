package bridges

import (
	"fmt"
	"github.com/kubemq-hub/builder/connector/bridges/source"
	"github.com/kubemq-hub/builder/connector/bridges/target"
	"github.com/kubemq-hub/builder/connector/common"
	"github.com/kubemq-hub/builder/pkg/utils"
	"github.com/kubemq-hub/builder/survey"
)

type Binding struct {
	Name              string            `json:"name"`
	Sources           *source.Source    `json:"sources"`
	Targets           *target.Target    `json:"targets"`
	Properties        map[string]string `json:"properties"`
	SourcesSpec       string            `json:"-" yaml:"-"`
	TargetsSpec       string            `json:"-" yaml:"-"`
	PropertiesSpec    string            `json:"-" yaml:"-"`
	addressOptions    []string
	takenSourceNames  []string
	takenTargetsNames []string
	takenBindingNames []string
	defaultName       string
	isEditMode        bool
	wasEdited         bool
}

func NewBinding(defaultName string) *Binding {
	return &Binding{
		defaultName: defaultName,
	}
}
func (b *Binding) Clone() *Binding {
	newBnd := &Binding{
		Name:              b.Name,
		Sources:           b.Sources.Clone(),
		Targets:           b.Targets.Clone(),
		Properties:        map[string]string{},
		SourcesSpec:       b.SourcesSpec,
		TargetsSpec:       b.TargetsSpec,
		PropertiesSpec:    b.PropertiesSpec,
		addressOptions:    b.addressOptions,
		takenSourceNames:  b.takenSourceNames,
		takenTargetsNames: b.takenTargetsNames,
		takenBindingNames: b.takenBindingNames,
		defaultName:       b.Name,
	}
	for key, val := range b.Properties {
		newBnd.Properties[key] = val
	}

	return newBnd
}
func (b *Binding) SetAddress(value []string) *Binding {
	b.addressOptions = value
	return b
}
func (b *Binding) SetEditMode(value bool) *Binding {
	b.isEditMode = value
	return b
}
func (b *Binding) SetDefaultName(value string) *Binding {
	b.defaultName = value
	return b
}
func (b *Binding) SetTakenSourceNames(value []string) *Binding {
	b.takenSourceNames = value
	return b
}
func (b *Binding) SetTakenTargetsNames(value []string) *Binding {
	b.takenTargetsNames = value
	return b
}
func (b *Binding) SetTakenBindingNames(value []string) *Binding {
	b.takenBindingNames = value
	return b
}


func (b *Binding) setSource() error {
	if !b.isEditMode {
		utils.Println(promptSourceStart)
		b.Sources = source.NewSource(fmt.Sprintf("%s-source", b.defaultName))
	}

	var err error
		if b.Sources, err = b.Sources.
			SetAddress(b.addressOptions).
			SetIsEdit(b.isEditMode).
			SetTakenNames(b.takenSourceNames).
			Render(); err != nil {
			return err
		}
	return nil
}
func (b *Binding) setTarget() error {
	if !b.isEditMode {
		utils.Println(promptTargetStart)
		b.Targets = target.NewTarget(fmt.Sprintf("%s-target", b.defaultName))
	}
	var err error
		if b.Targets, err = b.Targets.
			SetAddress(b.addressOptions).
			SetIsEdit(b.isEditMode).
			SetTakenNames(b.takenSourceNames).
			Render(); err != nil {
			return err
		}
	return nil
}

func (b *Binding) setProperties() error {
	var err error
		p := common.NewProperties()
		if b.Properties, err = p.
			Render(); err != nil {
			return err
		}
	b.PropertiesSpec = p.ColoredYaml()
	return nil
}
func (b *Binding) showConfiguration() error {
	utils.Println(promptShowBinding, b.Name)
	utils.Println(b.ColoredYaml())
	return nil
}
func (b *Binding) setName() error {
	var err error
	if b.Name, err = NewName(b.defaultName).
		SetTakenNames(b.takenBindingNames).
		Render(); err != nil {
		return err
	}
	b.wasEdited = true
	return nil
}
func (b *Binding) add() (*Binding, error) {
	if err := b.setName(); err != nil {
		return nil, err
	}

	if err := b.setSource(); err != nil {
		return nil, err
	}

	if err := b.setTarget(); err != nil {
		return nil, err
	}
	utils.Println(promptBindingComplete)
	if err := b.setProperties(); err != nil {
		return nil, err
	}
	return b, nil
}

func (b *Binding) edit() (*Binding, error) {
	var result *Binding
	edited := b.Clone().
		SetEditMode(true)

	form := survey.NewForm(fmt.Sprintf("Select Edit %s Binding Option:", edited.Name))

	ftName := new(string)
	*ftName = fmt.Sprintf("<n> Edit Binding's Name (%s)", edited.Name)
	form.AddItem(ftName, func() error {
		if err := edited.setName(); err != nil {
			return err
		}
		*ftName = fmt.Sprintf("<n> Edit Binding's Name (%s)", edited.Name)
		return nil
	})

	ftSource := new(string)
	*ftSource = fmt.Sprintf("<s> Edit Binding's Source (%s)", edited.Source.Kind)
	form.AddItem(ftSource, func() error {
		var err error
		if edited.Source, err = edited.editSource(); err != nil {
			return err
		}
		*ftSource = fmt.Sprintf("<s> Edit Binding's Source (%s)", edited.Source.Kind)
		return nil
	})

	ftTarget := new(string)
	*ftTarget = fmt.Sprintf("<t> Edit Binding's Target (%s)", edited.Target.Kind)
	form.AddItem(ftTarget, func() error {
		var err error
		if edited.Target, err = edited.editTarget(); err != nil {
			return err
		}
		*ftTarget = fmt.Sprintf("<t> Edit Binding's Target (%s)", edited.Target.Kind)
		return nil
	})

	form.AddItem("<m> Edit Binding's Middlewares", edited.setProperties)

	form.AddItem("<c> Show Binding Configuration", edited.showConfiguration)

	form.SetOnSaveFn(func() error {
		if err := edited.Validate(); err != nil {
			return err
		}
		result = edited
		return nil
	})

	form.SetOnCancelFn(func() error {
		result = b
		return nil
	})
	if err := form.Render(); err != nil {
		return nil, err
	}

	return result, nil

}
func (b *Binding) Render() (*Binding, error) {
	if b.isEditMode {
		return b.edit()
	}
	return b.add()
}

func (b *Binding) ColoredYaml() string {
	b.SourcesSpec = b.Sources.ColoredYaml()
	b.TargetsSpec = b.Targets.ColoredYaml()
	b.PropertiesSpec = utils.MapToYaml(b.Properties)
	tpl := utils.NewTemplate(bindingTemplate, b)
	bnd, err := tpl.Get()
	if err != nil {
		return fmt.Sprintf("error rendring binding spec,%s", err.Error())
	}
	return string(bnd)
}
