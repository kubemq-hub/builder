package cluster

import (
	"fmt"
	"github.com/kubemq-hub/builder/survey"
	"strings"
)

const nodeSelectorTml = `
<red>nodeSelectors:{{ range $key, $value := . }}
  {{ $key -}}: {{ $value -}}{{ end }}</>
`

type NodeSelector struct {
	values map[string]string
}

func NewNodeSelector() *NodeSelector {
	return &NodeSelector{}
}
func (n *NodeSelector) Clone() *NodeSelector {
	newValues := map[string]string{}
	for key, val := range n.values {
		newValues[key] = val
	}
	return &NodeSelector{
		values: newValues,
	}
}
func (n *NodeSelector) Validate() error {
	return nil
}
func (n *NodeSelector) Render(values map[string]string) (map[string]string, error) {
	if values != nil {
		n.values = values
	} else {
		n.values = map[string]string{}
	}

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
func (n *NodeSelector) ColoredYaml() (string, error) {
	t := NewTemplate(nodeSelectorTml, n)
	b, err := t.Get()
	if err != nil {
		return fmt.Sprintf("error rendring node selectors spec,%s", err.Error()), nil
	}
	return string(b), nil
}

var _ Validator = NewNodeSelector()
