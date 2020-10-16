package common

import (
	"fmt"
	"github.com/kubemq-hub/builder/survey"
	"gopkg.in/yaml.v2"
)

type Bindings struct {
	Bindings          []*Binding `json:"bindings"`
	manifest          *Manifest
	loadedOptions     DefaultOptions
	takenBindingNames []string
}

func NewBindings() *Bindings {
	return &Bindings{}
}

func (b *Bindings) SetManifest(value *Manifest) *Bindings {
	b.manifest = value
	return b
}
func (b *Bindings) SetOptions(value DefaultOptions) *Bindings {
	b.loadedOptions = value
	return b
}
func (b *Bindings) askAddBinding() (bool, error) {
	val := false
	err := survey.NewBool().
		SetKind("bool").
		SetName("add-binding").
		SetMessage("Would you like to add another binding").
		SetDefault("false").
		SetHelp("Add new binding").
		SetRequired(true).
		Render(&val)
	if err != nil {
		return false, err
	}
	return val, nil
}

func (b *Bindings) addBinding() error {
	if bnd, err := NewBinding().
		SetDefaultOptions(b.loadedOptions).
		SetSourcesList(b.manifest.Sources).
		SetTargetsList(b.manifest.Targets).
		SetTakenBindingNames(b.takenBindingNames).
		Render(); err != nil {
		return err
	} else {
		b.Bindings = append(b.Bindings, bnd)
		b.takenBindingNames = append(b.takenBindingNames, bnd.Name)
	}
	return nil
}

func (b *Bindings) Render() ([]byte, error) {
	if b.manifest == nil {
		return nil, fmt.Errorf("inavlid manifest")
	}
	err := b.addBinding()
	if err != nil {
		return nil, err
	}
	for {
		addMore, err := b.askAddBinding()
		if err != nil {
			return nil, err
		}
		if addMore {
			err := b.addBinding()
			if err != nil {
				return nil, err
			}
		} else {
			goto done
		}
	}
done:
	return yaml.Marshal(b)
}

func (b *Bindings) Yaml() ([]byte, error) {
	return yaml.Marshal(b)
}
