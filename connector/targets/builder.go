package targets

import (
	"fmt"
	"github.com/ghodss/yaml"
	"github.com/kubemq-hub/builder/connector/pkg/builder"
)

type Target struct {
	manifestData   []byte
	defaultOptions builder.DefaultOptions
}

func NewSource() *Target {
	return &Target{}
}

func (s *Target) SetManifest(value []byte) *Target {
	s.manifestData = value
	return s
}

func (s *Target) SetDefault(value builder.DefaultOptions) *Target {
	s.defaultOptions = value
	return s
}

func (s *Target) Render() ([]byte, error) {
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
