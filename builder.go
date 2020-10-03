package main

import (
	"github.com/kubemq-hub/builder/common"
	"github.com/kubemq-hub/builder/survey"
)

type Builder struct {
	kind           string
	defaultOptions common.DefaultOptions
}

func NewBuilder() *Builder {
	return &Builder{}
}

func (b *Builder) SetDefaultOptions(value common.DefaultOptions) *Builder {
	b.defaultOptions = value
	return b
}
func (b *Builder) askType() error {
	err := survey.NewString().
		SetKind("string").
		SetName("connector type").
		SetMessage("Choose connector type").
		SetOptions([]string{"KubeMQ Bridges", "KubeMQ Targets", "KubeMQ Sources"}).
		SetDefault("KubeMQ Bridges").
		SetHelp("Set connector type").
		SetRequired(true).
		Render(&b.kind)
	if err != nil {
		return err
	}
	return nil
}
