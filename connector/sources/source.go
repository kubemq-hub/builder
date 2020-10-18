package sources

import (
	"fmt"
	"github.com/kubemq-hub/builder/connector/common"
	"io/ioutil"
)

type Source struct {
	manifestData   []byte
	defaultOptions common.DefaultOptions
	defaultName    string
	bindingsList   []*common.Binding
}

func NewSource(defaultName string) *Source {
	return &Source{
		defaultName: defaultName,
	}
}

func (s *Source) SetManifest(value []byte) *Source {
	s.manifestData = value
	return s
}
func (s *Source) SetManifestFile(filename string) *Source {
	s.manifestData, _ = ioutil.ReadFile(filename)
	return s
}
func (s *Source) SetBindings(value []*common.Binding) *Source {
	s.bindingsList = value
	return s
}

func (s *Source) SetDefaultOptions(value common.DefaultOptions) *Source {
	s.defaultOptions = value
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
	if b, err := common.NewBindings(s.defaultName).
		SetManifest(m).
		SetOptions(s.defaultOptions).
		SetBindings(s.bindingsList).
		Render(); err != nil {
		return nil, err
	} else {
		return b, nil
	}
}
