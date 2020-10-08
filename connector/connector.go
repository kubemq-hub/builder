package connector

import (
	"github.com/kubemq-hub/builder/common"
	"github.com/kubemq-hub/builder/survey"
)

type Connector struct {
	Name            string `json:"name"`
	Namespace       string `json:"namespace"`
	Type            string `json:"type"`
	Replicas        int    `json:"replicas"`
	Config          string `json:"config"`
	NodePort        int    `json:"node_port"`
	ServiceType     string `json:"service_type"`
	Image           string `json:"image"`
	defaultOptions  common.DefaultOptions
	targetManifest  []byte
	sourcesManifest []byte
}

func NewConnector() *Connector {
	return &Connector{}
}
func (c *Connector) SetDefaultOptions(value common.DefaultOptions) *Connector {
	c.defaultOptions = value
	return c
}
func (c *Connector) SetTargetsManifest(value []byte) *Connector {
	c.targetManifest = value
	return c
}
func (c *Connector) SetSourcesManifest(value []byte) *Connector {
	c.sourcesManifest = value
	return c
}

func (c *Connector) askType() error {
	err := survey.NewString().
		SetKind("string").
		SetName("connector type").
		SetMessage("Choose Connector type").
		SetOptions([]string{"KubeMQ Bridges", "KubeMQ Targets", "KubeMQ Sources"}).
		SetDefault("KubeMQ Bridges").
		SetHelp("Set Connector type").
		SetRequired(true).
		Render(&c.Type)
	if err != nil {
		return err
	}
	return nil
}
func (c *Connector) askImage() error {
	err := survey.NewString().
		SetKind("string").
		SetName("connector image").
		SetMessage("Set Connector image").
		SetDefault("latest").
		SetHelp("Set Connector image").
		SetRequired(false).
		Render(&c.Image)
	if err != nil {
		return err
	}
	if c.Image == "latest" {
		c.Image = ""
	}
	return nil
}
func (c *Connector) askName() error {
	if name, err := NewName().
		Render(); err != nil {
		return err
	} else {
		c.Name = name.Name
	}
	return nil
}
func (c *Connector) askService() error {
	err := survey.NewString().
		SetKind("string").
		SetName("service-type").
		SetMessage("Set Connector service type").
		SetDefault("ClusterIP").
		SetOptions([]string{"ClusterIP", "NodePort", "LoadBalancer"}).
		SetHelp("Sets Connector service type").
		SetRequired(true).
		Render(&c.ServiceType)
	if err != nil {
		return err
	}
	if c.ServiceType != "NodePort" {
		return nil
	}
	err = survey.NewInt().
		SetKind("int").
		SetName("node-port").
		SetMessage("Set Connector service NodePort value").
		SetDefault("30000").
		SetHelp("Set Connector service NodePort value").
		SetRequired(false).
		SetRange(30000, 32767).
		Render(&c.NodePort)
	if err != nil {
		return err
	}

	return nil
}
func (c *Connector) askNamespace() error {
	if n, err := NewNamespace().
		SetNamespaces(c.defaultOptions["namespaces"]).
		Render(); err != nil {
		return err
	} else {
		c.Namespace = n.Namespace
	}
	return nil
}
func (c *Connector) askReplicas() error {
	var err error
	if c.Replicas, err = NewReplicas().
		Render(); err != nil {
		return err
	}
	return nil
}

func (c *Connector) Render() (*Connector, error) {
	if err := c.askType(); err != nil {
		return nil, err
	}
	if err := c.askName(); err != nil {
		return nil, err
	}
	if err := c.askNamespace(); err != nil {
		return nil, err
	}
	if err := c.askReplicas(); err != nil {
		return nil, err
	}
	if err := c.askService(); err != nil {
		return nil, err
	}
	if err := c.askImage(); err != nil {
		return nil, err
	}

	switch c.Type {
	case "KubeMQ Bridges":
		cfg, err := NewBridge().
			SetClusterAddress(c.defaultOptions["kubemq-address"]).
			Render()
		if err != nil {
			return nil, err
		}
		c.Config = string(cfg)
		c.Type = "bridges"
	case "KubeMQ Targets":
		cfg, err := NewTarget().
			SetManifest(c.targetManifest).
			SetDefaultOptions(c.defaultOptions).
			Render()
		if err != nil {
			return nil, err
		}
		c.Config = string(cfg)
		c.Type = "targets"
	case "KubeMQ Sources":
		cfg, err := NewSource().
			SetManifest(c.sourcesManifest).
			SetDefaultOptions(c.defaultOptions).
			Render()
		if err != nil {
			return nil, err
		}
		c.Config = string(cfg)
		c.Type = "source"
	}

	return c, nil
}
