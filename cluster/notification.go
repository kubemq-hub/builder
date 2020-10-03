package cluster

import (
	"github.com/kubemq-hub/builder/survey"
)

type Notification struct {
	Prefix  string `json:"prefix"`
	Enabled bool   `json:"enabled"`
}

func NewNotification() *Notification {
	return &Notification{}
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

var _ Validator = NewNotification()
