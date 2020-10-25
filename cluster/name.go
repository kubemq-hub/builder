package cluster

import (
	"fmt"
	"github.com/kubemq-hub/builder/survey"
)

type Name struct {
	Name       string
	takenNames []string
}

func NewName() *Name {
	return &Name{}
}
func (n *Name) Clone() *Name {
	return &Name{
		Name:       n.Name,
		takenNames: n.takenNames,
	}
}
func (n *Name) SetTakenNames(value []string) *Name {
	n.takenNames = value
	return n
}

func (n *Name) Validate() error {
	if n.Name == "" {
		return fmt.Errorf("name cannot be empty")
	}
	return nil
}
func (n *Name) checkNonEmptyName(val interface{}) error {
	str, _ := val.(string)
	if str == "" {
		return fmt.Errorf("name cannot be empty")
	}
	return nil
}
func (n *Name) Render(defaultName string) (*Name, error) {
	err := survey.NewString().
		SetKind("string").
		SetName("name").
		SetMessage("Set Cluster name").
		SetDefault(defaultName).
		SetHelp("Set cluster name").
		SetRequired(true).
		SetInvalidOptions(n.takenNames).
		SetInvalidOptionsMessage("Cluster name must be unique within the same namespace").
		SetValidator(n.checkNonEmptyName).
		Render(&n.Name)
	if err != nil {
		return nil, err
	}
	return n, nil
}

var _ Validator = NewName()
