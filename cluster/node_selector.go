package cluster

import (
	"github.com/kubemq-hub/builder/survey"
	"strings"
)

type NodeSelector struct {
	values map[string]string
}

func NewNodeSelector() *NodeSelector {
	return &NodeSelector{}
}

func (n *NodeSelector) Validate() error {
	return nil
}
func (n *NodeSelector) Render() (map[string]string, error) {
	n.values = map[string]string{}
	val := ""
	err := survey.NewString().
		SetKind("string").
		SetName("node-selector").
		SetMessage("Set node selectors (key1=value1;key2=value2;...)").
		SetDefault("").
		SetHelp("Set node selectors (key1=value1;key2=value2;...)").
		SetRequired(false).
		Render(&val)
	if err != nil {
		return nil, err
	}
	for _, str := range strings.Split(val, ";") {
		keyValue := strings.Split(str, "=")
		if len(keyValue) == 2 {
			n.values[keyValue[0]] = keyValue[1]
		}
	}
	return n.values, nil
}

var _ Validator = NewNodeSelector()
