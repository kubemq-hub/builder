package sources

import (
	"fmt"
	"github.com/kubemq-hub/builder/connector/common"
	"io/ioutil"
)

type Source struct {
	manifestData  []byte
	loadedOptions common.DefaultOptions
	defaultName   string
	bindings      []*common.Binding
}

func NewSource(defaultName string, bindings []*common.Binding, loadedOptions common.DefaultOptions, manifestData []byte) *Source {
	return &Source{
		manifestData:  manifestData,
		loadedOptions: loadedOptions,
		defaultName:   defaultName,
		bindings:      bindings,
	}
}

func (s *Source) SetManifestFile(filename string) *Source {
	s.manifestData, _ = ioutil.ReadFile(filename)
	return s
}

func (s *Source) Render() ([]byte, error) {
	if s.manifestData == nil {
		return nil, fmt.Errorf("invalid manifest")
	}
	m, err := common.LoadManifest(s.manifestData)
	if err != nil {
		return nil, err
	}
	if m.Schema != "sources" {
		return nil, fmt.Errorf("invalid scheme")
	}
	if b, err := common.NewBindings(s.defaultName, s.bindings, "sources", s.loadedOptions, m).
		Render(); err != nil {
		return nil, err
	} else {
		return b, nil
	}
}
