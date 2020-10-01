package survey

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
)

type KindMeta struct {
	Kind    string `json:"kind"`
	askOpts []survey.AskOpt
}

func NewKindMeta() *KindMeta {
	return &KindMeta{}
}

func (k *KindMeta) SetKind(value string) *KindMeta {
	k.Kind = value
	return k
}
func (k *KindMeta) complete() error {
	switch k.Kind {
	case "string":
	case "int":
	case "confirm":
	default:
		return fmt.Errorf("invalid kind")
	}
	return nil
}
