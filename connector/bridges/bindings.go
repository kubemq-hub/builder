package bridges

import (
	"fmt"
	"github.com/ghodss/yaml"
	"github.com/kubemq-hub/builder/connector/common"
	"github.com/kubemq-hub/builder/pkg/utils"
	"github.com/kubemq-hub/builder/survey"
)

type Bindings struct {
	Bindings          []*Binding `json:"bindings"`
	defaultOptions    common.DefaultOptions
	takenBindingNames []string
	takenSourceNames  []string
	takenTargetNames  []string
	addressOptions    []string
	defaultName       string
}

func NewBindings(defaultName string) *Bindings {
	return &Bindings{
		defaultName: defaultName,
	}
}

func (b *Bindings) SetDefaultOptions(value common.DefaultOptions) *Bindings {
	b.defaultOptions = value
	return b
}

func (b *Bindings) askAddBinding() (bool, error) {
	val := false
	err := survey.NewBool().
		SetKind("bool").
		SetName("add-binding").
		SetMessage("Would you like to add another bindings").
		SetDefault("false").
		SetHelp("Add new bindings bridge").
		SetRequired(true).
		Render(&val)
	if err != nil {
		return false, err
	}
	return val, nil
}
func (b *Bindings) confirmBinding(bnd *Binding) bool {
	utils.Println(fmt.Sprintf(promptBindingConfirm, bnd.String()))
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
func (b *Bindings) addBinding() error {

	bnd := NewBinding(fmt.Sprintf("%s-binding-%d", b.defaultName, len(b.Bindings)+1))
	var err error
	if bnd, err = bnd.
		SetAddress(b.addressOptions).
		SetTakenBindingNames(b.takenBindingNames).
		SetTakenSourceNames(b.takenSourceNames).
		SetTakenTargetsNames(b.takenTargetNames).
		Render(); err != nil {
		return err
	}
	ok := b.confirmBinding(bnd)
	if ok {
		b.Bindings = append(b.Bindings, bnd)
		b.takenBindingNames = append(b.takenBindingNames, bnd.BindingName())
		b.takenSourceNames = append(b.takenSourceNames, bnd.SourceName())
		b.takenTargetNames = append(b.takenTargetNames, bnd.TargetName())

	}
	return nil
}
func (b *Bindings) askSelectBinding(op string) (*Binding, error) {
	var bindingList []string
	for _, bnd := range b.Bindings {
		bindingList = append(bindingList, bnd.Name)
	}
	bindingList = append(bindingList, "Cancel")
	val := ""
	err := survey.NewString().
		SetKind("string").
		SetName("select-binding").
		SetMessage(fmt.Sprintf("Select Binding name to %s", op)).
		SetDefault(bindingList[0]).
		SetHelp("Select Binding name to delete or Cancel ").
		SetRequired(true).
		SetOptions(bindingList).
		Render(&val)
	if err != nil {
		return nil, err
	}
	if val == "Cancel" {
		return nil, nil
	}
	for _, binding := range b.Bindings {
		if val == binding.Name {
			return binding, nil
		}
	}
	return nil, nil
}
func (b *Bindings) switchOrRemove(old, new *Binding) {
	var newBindingList []*Binding
	var newTakenBindingNames []string
	var newTakenSourceNames []string
	var newTakenTargetNames []string

	for _, binding := range b.Bindings {
		if old.Name != binding.Name {
			newBindingList = append(newBindingList, binding)
			newTakenBindingNames = append(newTakenBindingNames, binding.Name)
			newTakenSourceNames = append(newTakenSourceNames, binding.Sources.Name)
			newTakenTargetNames = append(newTakenTargetNames, binding.Targets.Name)
		}
	}
	if new != nil {
		newBindingList = append(newBindingList, new)
		newTakenBindingNames = append(newTakenBindingNames, new.Name)
		newTakenBindingNames = append(newTakenBindingNames, new.Name)
		newTakenSourceNames = append(newTakenSourceNames, new.Sources.Name)
		newTakenTargetNames = append(newTakenTargetNames, new.Targets.Name)
	}
	b.Bindings = newBindingList
	b.takenBindingNames = newTakenBindingNames
	b.takenSourceNames = newTakenSourceNames
	b.takenTargetNames = newTakenTargetNames

}
func (b *Bindings) editBinding() error {
	bnd, err := b.askSelectBinding("edit")
	if err != nil {
		return err
	}

	if bnd == nil {
		utils.Println(promptBindingEditCanceled)
		return nil
	}

	edited := bnd.Clone()
	if edited, err = edited.
		SetEditMode(true).
		SetAddress(b.addressOptions).
		SetTakenBindingNames(b.takenBindingNames).
		SetTakenSourceNames(b.takenSourceNames).
		SetTakenTargetsNames(b.takenTargetNames).
		Render(); err != nil {
		return err
	}
	ok := b.confirmBinding(edited)
	if ok {
		b.switchOrRemove(bnd, edited)
		utils.Println(promptBindingEditedConfirmation, bnd.Name)

	} else {
		utils.Println(promptBindingEditedNoSave, bnd.Name)
	}

	return nil
}
func (b *Bindings) deleteBinding() error {
	bnd, err := b.askSelectBinding("delete")
	if err != nil {
		return err
	}
	if bnd == nil {
		utils.Println(promptBindingDeleteCanceled)
		return nil
	}
	b.switchOrRemove(bnd, nil)
	utils.Println(promptBindingDeleteConfirmation, bnd.Name)
	return nil
}
func (b *Bindings) askMenu() error {
	utils.Println(promptBindingStartMenu)
	ops := []string{
		"Add Binding",
		"Edit Binding",
		"Delete Binding",
		"Done",
	}
	for {
		val := ""
		err := survey.NewString().
			SetKind("string").
			SetName("select-operation").
			SetMessage("Select Binding operation").
			SetDefault(ops[0]).
			SetHelp("Select Binding operation").
			SetRequired(true).
			SetOptions(ops).
			Render(&val)
		if err != nil {
			return err
		}
		switch val {
		case "Add Binding":
			if err := b.addBinding(); err != nil {
				return err
			}
		case "Edit Binding":
			if err := b.editBinding(); err != nil {
				return err
			}
		case "Delete Binding":
			if err := b.deleteBinding(); err != nil {
				return err
			}
		default:
			return nil
		}
	}
}
func (b *Bindings) Render() ([]byte, error) {
	if err := b.askMenu(); err != nil {
		return nil, err
	}

	if len(b.Bindings) == 0 {
		return nil, fmt.Errorf("at least one binding must be configured")
	}
	return yaml.Marshal(b)
}

func (b *Bindings) Marshal() ([]byte, error) {
	return yaml.Marshal(b)
}
func (b *Bindings) Unmarshal(data []byte) *Bindings {
	bnd := &Bindings{}
	err := yaml.Unmarshal(data, bnd)
	if err != nil {
		return b
	}
	return bnd
}
