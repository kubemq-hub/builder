package cluster

import (
	"fmt"
	"github.com/kubemq-hub/builder/survey"
	"math"
)

const authorizationTml = `
<red>authorization:</>
  <yellow>enable:</> <white>{{ .Enable }}</>  
  <yellow>policy:</> |-<white>{{ .Policy | indent 4}}</>
  <yellow>url:</> <white>{{ .Url}}</>
  <yellow>autoReload:</> <white>{{ .AutoReload}}</>
`

type Authorization struct {
	Policy     string `json:"policy"`
	Url        string `json:"url"`
	AutoReload int    `json:"auto_reload"`
}

func NewAuthorization() *Authorization {
	return &Authorization{}
}
func (a *Authorization) Clone() *Authorization {
	return &Authorization{
		Policy:     a.Policy,
		Url:        a.Url,
		AutoReload: a.AutoReload,
	}
}
func (a *Authorization) askPolicy() error {
	err := survey.NewMultiline().
		SetKind("multiline").
		SetName("policy").
		SetMessage("Load Authorization policy data").
		SetDefault("").
		SetHelp("Load Authorization policy data").
		SetRequired(false).
		Render(&a.Policy)
	if err != nil {
		return err
	}
	return nil
}
func (a *Authorization) askUrl() error {
	err := survey.NewString().
		SetKind("string").
		SetName("url").
		SetMessage("Set URL policy loading address").
		SetDefault("").
		SetHelp("Set URL policy loading address").
		SetRequired(false).
		Render(&a.Url)
	if err != nil {
		return err
	}
	return nil
}
func (a *Authorization) askAutoReload() error {
	err := survey.NewInt().
		SetKind("int").
		SetName("auto-reload").
		SetMessage("Set automatic policy reload seconds (0 - no reload)").
		SetDefault("0").
		SetHelp("Set automatic policy reload seconds (0 - no reload)").
		SetRequired(false).
		SetRange(0, math.MaxInt32).
		Render(&a.AutoReload)
	if err != nil {
		return err
	}
	return nil
}

func (a *Authorization) Validate() error {
	if a.AutoReload < 0 {
		return fmt.Errorf("auto reload value cannot be less than 0")
	}
	return nil
}

func (a *Authorization) Render() (*Authorization, error) {
	if err := a.askPolicy(); err != nil {
		return nil, err
	}
	if err := a.askUrl(); err != nil {
		return nil, err
	}

	if err := a.askAutoReload(); err != nil {
		return nil, err
	}
	return a, nil
}
func (a *Authorization) ColoredYaml() (string, error) {
	t := NewTemplate(authorizationTml, a)
	b, err := t.Get()
	if err != nil {
		return fmt.Sprintf("error rendring authentication spec,%s", err.Error()), nil
	}
	return string(b), nil
}

var _ Validator = NewAuthorization()
