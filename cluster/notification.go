package cluster

import (
	"fmt"
	"github.com/kubemq-hub/builder/survey"
)

const notificationTmpl = `
<red>notification:</>
  <yellow>enabled:</> <white>{{ .Enable }}</>
  <yellow>prefix:</> <white>{{ .Prefix}}</>
`

type Notification struct {
	Prefix  string `json:"prefix"`
	Enabled bool   `json:"enabled"`
}

func NewNotification() *Notification {
	return &Notification{}
}
func (n *Notification) Clone() *Notification {
	return &Notification{
		Prefix:  n.Prefix,
		Enabled: n.Enabled,
	}
}
func (n *Notification) askPrefix() error {
	err := survey.NewString().
		SetKind("string").
		SetName("prefix").
		SetMessage("Set notification reporting channel prefix").
		SetDefault("").
		SetHelp("Set notification reporting channel prefix").
		SetRequired(false).
		Render(&n.Prefix)
	if err != nil {
		return err
	}
	return nil
}

func (n *Notification) Validate() error {
	return nil
}
func (n *Notification) Render() (*Notification, error) {
	if err := n.askPrefix(); err != nil {
		return nil, err
	}
	if n.Prefix != "" {
		n.Enabled = true
	}
	return n, nil
}
func (n *Notification) ColoredYaml() (string, error) {
	t := NewTemplate(notificationTmpl, n)
	b, err := t.Get()
	if err != nil {
		return fmt.Sprintf("error rendring notification selectors spec,%s", err.Error()), nil
	}
	return string(b), nil
}

var _ Validator = NewNotification()
