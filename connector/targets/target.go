package targets

import (
	"fmt"
	"github.com/kubemq-hub/builder/connector/common"
	"io/ioutil"
)

type Target struct {
	manifestData  []byte
	loadedOptions common.DefaultOptions
	defaultName   string
	bindings      []*common.Binding
}

func NewTarget(defaultName string, bindings []*common.Binding, loadedOptions common.DefaultOptions, manifestData []byte) *Target {
	return &Target{
		manifestData:  manifestData,
		loadedOptions: loadedOptions,
		defaultName:   defaultName,
		bindings:      bindings,
	}
}

func (t *Target) SetManifestFile(filename string) *Target {
	t.manifestData, _ = ioutil.ReadFile(filename)
	return t
}

func (t *Target) Render() ([]byte, error) {
	if t.manifestData == nil {
		return nil, fmt.Errorf("invalid manifest")
	}
	m, err := common.LoadManifest(t.manifestData)
	if err != nil {
		return nil, err
	}
	if m.Schema != "targets" {
		return nil, fmt.Errorf("invalid scheme, %s", m.Schema)
	}
	if b, err := common.NewBindings(t.defaultName, t.bindings, "targets", t.loadedOptions, m).
		Render(); err != nil {
		return nil, err
	} else {
		return b, nil
	}
}
