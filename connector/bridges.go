package connector

import (
	"fmt"
	"github.com/ghodss/yaml"
	"github.com/kubemq-hub/builder/connector/bridges/binding"
	"github.com/kubemq-hub/builder/pkg/utils"
	"github.com/kubemq-hub/builder/survey"
)

type Bridges struct {
	Bindings          []*binding.Binding `json:"bindings"`
	takenBindingNames []string
	takenSourceNames  []string
	takenTargetNames  []string
	addressOptions    []string
	defaultName       string
}

func NewBridges(defaultName string) *Bridges {
	return &Bridges{
		defaultName: defaultName,
	}
}
func (b *Bridges) SetClusterAddress(value []string) *Bridges {
	b.addressOptions = value
	return b
}

func (b *Bridges) askAddBinding() (bool, error) {
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
func (b *Bridges) confirmBinding(bnd *binding.Binding) bool {
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
	if !val {
		utils.Println(promptBindingReconfigure)
	}
	return val
}
func (b *Bridges) addBinding() error {
	for {
		bnd := binding.NewBinding(fmt.Sprintf("%s-binding-%d", b.defaultName, len(b.Bindings)+1))
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
			break
		}
	}

	return nil
}
func (b *Bridges) Render() ([]byte, error) {

	err := b.addBinding()
	if err != nil {
		return nil, err
	}
	for {
		addMore, err := b.askAddBinding()
		if err != nil {
			return yaml.Marshal(b)
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

func (b *Bridges) Yaml() ([]byte, error) {
	return yaml.Marshal(b)
}
