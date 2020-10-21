package connector

import (
	"fmt"
	"github.com/kubemq-hub/builder/connector/bridges"
	"github.com/kubemq-hub/builder/connector/common"
	"github.com/kubemq-hub/builder/connector/sources"
	"github.com/kubemq-hub/builder/connector/targets"
	"github.com/kubemq-hub/builder/pkg/utils"
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
	isEditMode      bool
}

func NewConnector() *Connector {
	return &Connector{}
}
func (c *Connector) Clone() *Connector {
	return &Connector{
		Name:            c.Name,
		Namespace:       c.Namespace,
		Type:            c.Type,
		Replicas:        c.Replicas,
		Config:          c.Config,
		NodePort:        c.NodePort,
		ServiceType:     c.ServiceType,
		Image:           c.Image,
		defaultOptions:  c.defaultOptions,
		targetManifest:  c.targetManifest,
		sourcesManifest: c.sourcesManifest,
		isEditMode:      false,
	}
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
func (c *Connector) Key() string {
	return fmt.Sprintf("%s/%s", c.Namespace, c.Name)
}
func (c *Connector) SetEditMode() *Connector {
	c.isEditMode = true
	return c
}
func (c *Connector) Confirm() bool {
	utils.Println(fmt.Sprintf(promptConnectorConfirm, c.ColoredYaml()))
	val := true
	err := survey.NewBool().
		SetKind("bool").
		SetName("confirm-connection").
		SetMessage("Would you like save this configuration").
		SetDefault("true").
		SetRequired(true).
		Render(&val)
	if err != nil {
		return false
	}
	if !val {
		utils.Println(promptConnectorReconfigure)
	}
	return val
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
func (c *Connector) askName(defaultName string) error {
	if name, err := NewName().
		Render(defaultName); err != nil {
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
		SetHelp("Set Connector service type").
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
func (c *Connector) ColoredYaml() string {
	t := utils.NewTemplate(connectorTemplate, c)
	b, err := t.Get()
	if err != nil {
		return fmt.Sprintf("error rendring source  spec,%s", err.Error())
	}
	return string(b)
}

func (c *Connector) Render() (*Connector, error) {
	if c.isEditMode {
		return c.edit()
	} else {
		return c.add()
	}
}
func (c *Connector) edit() (*Connector, error) {
	switch c.Type {
	case "bridges":
		bindings, err := bridges.Unmarshal([]byte(c.Config))
		if err != nil {
			return nil, err
		}
		bindings.SetDefaultOptions(c.defaultOptions)
		bindings.SetDefaultName(c.Name)
		if newBindings, err := bindings.Render(); err != nil {
			return nil, err
		} else {
			c.Config = string(newBindings)
		}
	case "targets", "sources":
		bindings, err := common.Unmarshal([]byte(c.Config))
		if err != nil {
			return nil, err
		}
		bindings.SetDefaultOptions(c.defaultOptions)
		bindings.SetDefaultName(c.Name)
		if newBindings, err := bindings.Render(); err != nil {
			return nil, err
		} else {
			c.Config = string(newBindings)
		}
	}
	utils.Println(promptConnectorContinue)
	if err := c.askReplicas(); err != nil {
		return nil, err
	}
	if err := c.askService(); err != nil {
		return nil, err
	}
	if err := c.askImage(); err != nil {
		return nil, err
	}
	return c, nil
}
func (c *Connector) add() (*Connector, error) {
	utils.Println(promptConnectorStart)
	menu := survey.NewMenu("Select Connector type").
		AddItem("KubeMQ Bridges", c.addBridges).
		AddItem("KubeMQ Targets", c.addTargets).
		AddItem("KubeMQ Sources", c.addSources).
		SetDisableLoop(true).SetBackOption(true).
		SetErrorHandler(survey.MenuShowErrorFn)
	if err := menu.Render(); err != nil {
		return nil, err
	}

	utils.Println(promptConnectorContinue)
	if err := c.askReplicas(); err != nil {
		return nil, err
	}
	if err := c.askService(); err != nil {
		return nil, err
	}
	if err := c.askImage(); err != nil {
		return nil, err
	}

	return c, nil
}

func (c *Connector) addBridges() error {
	if err := c.askName("kubemq-bridges"); err != nil {
		return err
	}
	if err := c.askNamespace(); err != nil {
		return err
	}
	utils.Println(promptBindingStart, c.Name)
	cfg, err := bridges.NewBridges(c.Name).
		SetDefaultOptions(c.defaultOptions).
		Render()
	if err != nil {
		return err
	}
	c.Config = string(cfg)
	c.Type = "bridges"
	return nil
}
func (c *Connector) addTargets() error {
	if err := c.askName("kubemq-targets"); err != nil {
		return err
	}
	if err := c.askNamespace(); err != nil {
		return err
	}
	utils.Println(promptBindingStart, c.Name)
	cfg, err := targets.NewTarget(c.Name).
		SetManifest(c.targetManifest).
		SetDefaultOptions(c.defaultOptions).
		Render()
	if err != nil {
		return err
	}
	c.Config = string(cfg)
	c.Type = "targets"
	return nil
}
func (c *Connector) addSources() error {
	if err := c.askName("kubemq-sources"); err != nil {
		return err
	}
	if err := c.askNamespace(); err != nil {
		return err
	}
	utils.Println(promptBindingStart, c.Name)
	cfg, err := sources.NewSource(c.Name).
		SetManifest(c.sourcesManifest).
		SetDefaultOptions(c.defaultOptions).
		Render()
	if err != nil {
		return err
	}
	c.Config = string(cfg)
	c.Type = "source"
	return nil
}
