package source

import "github.com/kubemq-hub/builder/connector/pkg/survey"

type Name struct {
	takenNames []string
}

func NewName() *Name {
	return &Name{}
}
func (n *Name) SetTakenNames(value []string) *Name {
	n.takenNames = value
	return n
}
func (n *Name) Render() (string, error) {
	val := ""
	err := survey.NewInput().
		SetKind("string").
		SetName("name").
		SetMessage("Set Source name").
		SetDefault("").
		SetHelp("Sets source name entry").
		SetRequired(true).
		SetInvalidOptions(n.takenNames).
		SetInvalidOptionsMessage("source name must be unique").
		Render(&val)
	if err != nil {
		return "", err
	}
	return val, nil
}
