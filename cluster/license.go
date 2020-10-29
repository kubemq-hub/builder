package cluster

import (
	"fmt"
	"github.com/kubemq-hub/builder/survey"
)

const licenseTml = `
<red>license</>: <white>|-
{{ .Data | indent 2}}/>
`

type License struct {
	Data string
}

func NewLicense() *License {
	return &License{}
}
func (l *License) Clone() *License {
	return &License{
		Data: l.Data,
	}
}
func (l *License) Validate() error {
	if l.Data == "" {
		return fmt.Errorf("license data cannot be empty")
	}
	return nil
}
func (l *License) Render(current string) (string, error) {
	err := survey.NewMultiline().
		SetKind("multiline").
		SetName("license").
		SetMessage("Load license data").
		SetDefault(current).
		SetHelp("Set license data via editor").
		SetRequired(false).
		Render(&l.Data)
	if err != nil {
		return "", err
	}
	return l.Data, nil
}
func (l *License) ColoredYaml() (string, error) {
	t := NewTemplate(licenseTml, l)
	b, err := t.Get()
	if err != nil {
		return fmt.Sprintf("error rendring license spec,%s", err.Error()), nil
	}
	return string(b), nil
}

var _ Validator = NewLicense()
