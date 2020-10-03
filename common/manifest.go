package common

import (
	"encoding/json"
	"io/ioutil"
)

type Manifest struct {
	Schema  string       `json:"schema"`
	Version string       `json:"version"`
	Sources []*Connector `json:"sources"`
	Targets []*Connector `json:"targets"`
}

func NewManifest() *Manifest {
	return &Manifest{}
}
func LoadManifest(data []byte) (*Manifest, error) {
	m := &Manifest{}
	err := json.Unmarshal(data, m)
	if err != nil {
		return nil, err
	}
	return m, nil
}
func (m *Manifest) Save(filename string) error {
	b, err := json.Marshal(m)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filename, b, 0644)
}

func (m *Manifest) SetSchema(value string) *Manifest {
	m.Schema = value
	return m
}

func (m *Manifest) SetVersion(value string) *Manifest {
	m.Version = value
	return m
}
func (m *Manifest) SetSourceConnectors(value []*Connector) *Manifest {
	m.Sources = value
	return m
}
func (m *Manifest) SetTargetConnectors(value []*Connector) *Manifest {
	m.Targets = value
	return m
}
func (m *Manifest) AddConnector(value *Connector) *Manifest {
	m.Sources = append(m.Sources, value)
	return m
}
