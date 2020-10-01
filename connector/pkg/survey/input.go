package survey

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"strconv"
)

type Input struct {
	*KindMeta
	*ObjectMeta
	Options               []string
	InvalidOptions        []string
	InvalidOptionsMessage string
	Range                 bool
	Min                   int
	Max                   int
	askOpts               []survey.AskOpt
	validators            []func(val interface{}) error
}

func NewInput() *Input {
	return &Input{
		KindMeta:   NewKindMeta(),
		ObjectMeta: NewObjectMeta(),
		Options:    nil,
		Range:      false,
		Min:        0,
		Max:        0,
		askOpts:    nil,
	}
}

func (i *Input) NewKindMeta() *Input {
	i.KindMeta = NewKindMeta()
	return i
}
func (i *Input) NewObjectMeta() *Input {
	i.ObjectMeta = NewObjectMeta()
	return i
}
func (i *Input) SetKind(value string) *Input {
	i.KindMeta.SetKind(value)
	return i
}

func (i *Input) SetName(value string) *Input {
	i.ObjectMeta.SetName(value)
	return i
}

func (i *Input) SetMessage(value string) *Input {
	i.ObjectMeta.SetMessage(value)
	return i
}

func (i *Input) SetDefault(value string) *Input {
	i.ObjectMeta.SetDefault(value)
	return i
}

func (i *Input) SetHelp(value string) *Input {
	i.ObjectMeta.SetHelp(value)
	return i
}
func (i *Input) SetInvalidOptionsMessage(value string) *Input {
	i.InvalidOptionsMessage = value
	return i
}
func (i *Input) SetInvalidOptions(value []string) *Input {
	i.InvalidOptions = value
	return i
}
func (i *Input) SetRequired(value bool) *Input {
	i.ObjectMeta.SetRequired(value)
	return i
}

func (i *Input) SetOptions(value []string) *Input {
	i.Options = value
	return i
}
func (i *Input) SetRange(value bool) *Input {
	i.Range = value
	return i
}
func (i *Input) SetValidator(f func(val interface{}) error) *Input {
	i.validators = append(i.validators, f)
	return i
}

func (i *Input) SetMin(value int) *Input {
	i.Min = value
	return i
}

func (i *Input) SetMax(value int) *Input {
	i.Max = value
	return i
}
func (i *Input) checkValue(val interface{}) error {
	if str, ok := val.(string); ok {
		val, err := strconv.Atoi(str)
		if err != nil {
			return fmt.Errorf("invalid integer")
		}
		if i.Range {
			if val < i.Min {
				return fmt.Errorf("value cannot be lower than minimum %d", i.Min)
			}
			if val > i.Max {
				return fmt.Errorf("value cannot be higer than maximum %d", i.Max)
			}
		}

	}
	return nil
}
func (i *Input) invalidOptionValidator(val interface{}) error {
	if str, ok := val.(string); ok {
		for _, item := range i.InvalidOptions {
			if str == item {
				return fmt.Errorf("%s", i.InvalidOptionsMessage)
			}
		}
	}
	return nil
}
func (i *Input) Complete() error {
	if err := i.KindMeta.complete(); err != nil {
		return err
	}
	i.askOpts = append(i.askOpts, i.KindMeta.askOpts...)

	if err := i.ObjectMeta.complete(); err != nil {
		return err
	}
	i.askOpts = append(i.askOpts, i.ObjectMeta.askOpts...)
	switch i.Kind {
	case "string":

	case "int":
		i.askOpts = append(i.askOpts, survey.WithValidator(i.checkValue))
	}
	if i.InvalidOptionsMessage == "" {
		i.InvalidOptionsMessage = "invalid option,"
	}
	if len(i.InvalidOptions) > 0 {
		i.askOpts = append(i.askOpts, survey.WithValidator(i.invalidOptionValidator))
	}
	for _, validator := range i.validators {
		i.askOpts = append(i.askOpts, survey.WithValidator(validator))
	}
	return nil
}

func (i *Input) Render(target interface{}) error {
	if err := i.Complete(); err != nil {
		return err
	}
	if len(i.Options) == 0 {
		singleInput := &survey.Input{
			Renderer: survey.Renderer{},
			Message:  i.Message,
			Default:  i.Default,
			Help:     i.Help,
		}
		return survey.AskOne(singleInput, target, i.askOpts...)
	}
	selectInput := &survey.Select{
		Renderer: survey.Renderer{},
		Message:  i.Message,
		Options:  i.Options,
		Default:  i.Default,
		Help:     i.Help,
	}
	err := survey.AskOne(selectInput, target, i.askOpts...)
	if err != nil {
		return err
	}
	val, _ := target.(*string)
	if *val == "Other" {
		singleInput := &survey.Input{
			Renderer: survey.Renderer{},
			Message:  fmt.Sprintf("%s, Other", i.Message),
			Default:  "",
			Help:     i.Help,
		}
		return survey.AskOne(singleInput, target, i.askOpts...)
	}
	return nil
}
