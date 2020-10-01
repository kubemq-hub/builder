package survey

import "github.com/AlecAivazis/survey/v2"

type Confirm struct {
	*KindMeta
	*ObjectMeta
	askOpts []survey.AskOpt
}

func (c *Confirm) NewKindMeta() *Confirm {
	c.KindMeta = NewKindMeta()
	return c
}
func (c *Confirm) NewObjectMeta() *Confirm {
	c.ObjectMeta = NewObjectMeta()
	return c
}
func (c *Confirm) SetKind(value string) *Confirm {
	c.KindMeta.SetKind(value)
	return c
}

func (c *Confirm) SetName(value string) *Confirm {
	c.ObjectMeta.SetName(value)
	return c
}

func (c *Confirm) SetMessage(value string) *Confirm {
	c.ObjectMeta.SetMessage(value)
	return c
}

func (c *Confirm) SetDefault(value string) *Confirm {
	c.ObjectMeta.SetDefault(value)
	return c
}

func (c *Confirm) SetHelp(value string) *Confirm {
	c.ObjectMeta.SetHelp(value)
	return c
}
func (c *Confirm) SetRequired(value bool) *Confirm {
	c.ObjectMeta.SetRequired(value)
	return c
}

func (c *Confirm) Complete() error {
	return nil
}

func (c *Confirm) Render(target interface{}) error {
	if err := c.Complete(); err != nil {
		return err
	}
	defValue := false
	if c.Default == "true" {
		defValue = true
	}
	confirm := &survey.Confirm{
		Renderer: survey.Renderer{},
		Message:  c.Message,
		Default:  defValue,
		Help:     c.Help,
	}
	return survey.AskOne(confirm, target, c.askOpts...)
}

func NewConfirm() *Confirm {
	return &Confirm{
		KindMeta:   NewKindMeta(),
		ObjectMeta: NewObjectMeta(),
	}
}
