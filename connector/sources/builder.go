package sources

import (
	"fmt"
	"github.com/ghodss/yaml"
	"github.com/kubemq-hub/builder/connector/pkg/builder"
)

type Source struct {
	manifestData   []byte
	defaultOptions builder.DefaultOptions
}

func NewSource() *Source {
	return &Source{}
}

func (s *Source) SetManifest(value []byte) *Source {
	s.manifestData = value
	return s
}

func (s *Source) SetDefault(value builder.DefaultOptions) *Source {
	s.defaultOptions = value
	return s
}

func (s *Source) Render() ([]byte, error) {
	if s.manifestData == nil {
		return nil, fmt.Errorf("invalid manifest")
	}
	m, err := builder.LoadManifest(s.manifestData)
	if err != nil {
		return nil, err
	}
	if b, err := builder.NewBuilder().
		SetManifest(m).
		SetOptions(s.defaultOptions).
		Render(); err != nil {
		return nil, err
	} else {
		return yaml.Marshal(b)
	}
}
