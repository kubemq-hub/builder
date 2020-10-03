package cluster

import (
	"fmt"
	"github.com/kubemq-hub/builder/survey"
)

type License struct {
	data string
}

func NewLicense() *License {
	return &License{}
}
func (l *License) Validate() error {
	if l.data == "" {
		return fmt.Errorf("license data cannot be empty")
	}
	return nil
}
func (l *License) Render() (string, error) {
	err := survey.NewMultiline().
		SetKind("multiline").
		SetName("license").
		SetMessage("Load license data").
		SetDefault("").
		SetHelp("Sets license data via editor").
		SetRequired(false).
		Render(&l.data)
	if err != nil {
		return "", err
	}
	return l.data, nil
}

var _ Validator = NewLicense()
