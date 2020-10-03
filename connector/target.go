package connector

import (
	"fmt"
	"github.com/ghodss/yaml"
	"github.com/kubemq-hub/builder/common"
)

type Target struct {
	manifestData   []byte
	defaultOptions common.DefaultOptions
}

func NewTarget() *Target {
	return &Target{}
}

func (s *Target) SetManifest(value []byte) *Target {
	s.manifestData = value
	return s
}

func (s *Target) SetDefaultOptions(value common.DefaultOptions) *Target {
	s.defaultOptions = value
	return s
}

func (s *Target) Render() ([]byte, error) {
	if s.manifestData == nil {
		return nil, fmt.Errorf("invalid manifest")
	}
	m, err := common.LoadManifest(s.manifestData)
	if err != nil {
		return nil, err
	}
	if b, err := common.NewBuilder().
		SetManifest(m).
		SetOptions(s.defaultOptions).
		Render(); err != nil {
		return nil, err
	} else {
		return yaml.Marshal(b)
	}
}
