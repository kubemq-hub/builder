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
	handler         ConnectorsHandler
}

func NewConnector(handler ConnectorsHandler) *Connector {
	return &Connector{
		handler: handler,
	}
}
func (c *Connector) Clone(handler ConnectorsHandler) *Connector {
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
		handler:         handler,
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

func (c *Connector) Validate() error {
	if c.Name == "" {
		return fmt.Errorf("connector must have a name")
	}
	if c.Namespace == "" {
		return fmt.Errorf("connector must have a namespace")
	}
	switch c.Type {
	case "bridges":
		bindings, err := bridges.Unmarshal([]byte(c.Config))
		if err != nil {
			return err
		}
		if err := bindings.Validate(); err != nil {
			return err
		}
	case "targets", "sources":
		bindings, err := common.Unmarshal([]byte(c.Config))
		if err != nil {
			return err
		}
		if err := bindings.Validate(); err != nil {
			return err
		}
	default:
		return fmt.Errorf("no valid connector type, %s", c.Type)
	}
	if c.Replicas < 0 {
		return fmt.Errorf("conenctor replicas cannot be negative")
	}
	if c.ServiceType == "" {
		return fmt.Errorf("conenctor service type must have a value")
	}
	if c.NodePort < 0 {
		return fmt.Errorf("conenctor node-port cannot be negative")
	}
	return nil
}
func (c *Connector) askImage() error {
	err := survey.NewString().
		SetKind("string").
		SetName("connector image").
		SetMessage("Set Connector image").
		SetDefault(c.Image).
		SetHelp("Set Connector image").
		SetRequired(false).
		Render(&c.Image)
	if err != nil {
		return err
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
		Render(c.Replicas); err != nil {
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

func EditConnector(origin *Connector, sourceManifest, targetsManifests []byte, handler ConnectorsHandler) (*Connector, error) {
	var result *Connector
	cloned := origin.Clone(handler).
		SetTargetsManifest(targetsManifests).
		SetSourcesManifest(sourceManifest)
	ftReplicas := new(string)
	*ftReplicas = fmt.Sprintf("Edit Connector Replicas (%d)", cloned.Replicas)

	ftServiceType := new(string)
	*ftServiceType = fmt.Sprintf("Edit Connector Service Type (%s)", cloned.ServiceType)

	ftImage := new(string)
	*ftImage = fmt.Sprintf("Edit Connector Docker Image (%s)", cloned.Image)

	form := survey.NewForm("Select Edit Connector Option:").
		AddItem(fmt.Sprintf("Edit Connector Bindings (%s)", cloned.Type), func() error {
			switch cloned.Type {
			case "bridges":
				bindings, err := bridges.Unmarshal([]byte(cloned.Config))
				if err != nil {
					return err
				}
				bindings.SetDefaultOptions(cloned.defaultOptions)
				bindings.SetDefaultName(cloned.Name)
				if newBindings, err := bindings.Render(); err != nil {
					return err
				} else {
					cloned.Config = string(newBindings)
				}
			case "targets", "sources":
				bindings, err := common.Unmarshal([]byte(cloned.Config))
				if err != nil {
					return err
				}
				bindings.SetDefaultOptions(cloned.defaultOptions)
				bindings.SetDefaultName(cloned.Name)

				if newBindings, err := bindings.Render(); err != nil {
					return err
				} else {
					cloned.Config = string(newBindings)
				}
			}
			utils.Println(promptConnectorContinue)
			return nil
		}).
		AddItem(ftReplicas, func() error {
			if err := cloned.askReplicas(); err != nil {
				return err
			}
			*ftReplicas = fmt.Sprintf("Edit Connector Replicas (%d)", cloned.Replicas)
			return nil
		}).
		AddItem(ftServiceType, func() error {
			if err := cloned.askService(); err != nil {
				return err
			}
			*ftServiceType = fmt.Sprintf("Edit Connector Service Type (%s)", cloned.ServiceType)
			return nil
		}).
		AddItem(ftImage, func() error {
			if err := cloned.askImage(); err != nil {
				return err

			}
			*ftImage = fmt.Sprintf("Edit Connector Docker Image (%s)", cloned.Image)
			return nil
		}).
		AddItem("Show Connector Configuration", func() error {
			utils.Println(cloned.ColoredYaml())
			return nil
		})

	form.SetOnSaveFn(func() error {
		if !(origin.ColoredYaml() == cloned.ColoredYaml()) {
			if err := cloned.Validate(); err != nil {
				return err
			}
			err := cloned.handler.Edit(cloned)
			if err != nil {
				return err
			}
			result = cloned

		} else {
			result = origin
		}
		return nil
	})
	form.SetOnCancelFn(func() error {
		result = origin
		return nil
	})

	form.SetOnErrorFn(survey.FormShowErrorFn)
	if err := form.Render(); err != nil {
		return nil, err
	}
	return result, nil
}

func AddConnector(sourceManifest, targetsManifests []byte, handler ConnectorsHandler) (*Connector, error) {
	var result *Connector
	added := NewConnector(handler).
		SetTargetsManifest(targetsManifests).
		SetSourcesManifest(sourceManifest)

	added.Replicas = 1
	added.ServiceType = "ClusterIP"
	ftReplicas := new(string)
	*ftReplicas = fmt.Sprintf("Set Connector Replicas (%d)", added.Replicas)

	ftServiceType := new(string)
	*ftServiceType = fmt.Sprintf("Set Connector Service Type (%s)", added.ServiceType)

	ftImage := new(string)
	*ftImage = fmt.Sprintf("Set Connector Docker Image (%s)", added.Image)

	form := survey.NewForm("Select Add Connector Option:").
		AddItem("Set Connector Type and Bindings", func() error {
			utils.Println(promptConnectorStart)
			menu := survey.NewMenu("Select Connector type").
				AddItem("KubeMQ Bridges", added.addBridges).
				AddItem("KubeMQ Targets", added.addTargets).
				AddItem("KubeMQ Sources", added.addSources).
				SetDisableLoop(true).SetBackOption(true).
				SetErrorHandler(survey.MenuShowErrorFn)
			if err := menu.Render(); err != nil {
				return err
			}
			return nil
		}).
		AddItem(ftReplicas, func() error {
			if err := added.askReplicas(); err != nil {
				return err
			}
			*ftReplicas = fmt.Sprintf("Edit Connector Replicas (%d)", added.Replicas)
			return nil
		}).
		AddItem(ftServiceType, func() error {
			if err := added.askService(); err != nil {
				return err
			}
			*ftServiceType = fmt.Sprintf("Edit Connector Service Type (%s)", added.ServiceType)
			return nil
		}).
		AddItem(ftImage, func() error {
			if err := added.askImage(); err != nil {
				return err

			}
			*ftImage = fmt.Sprintf("Edit Connector Docker Image (%s)", added.Image)
			return nil
		}).
		AddItem("Show Connector Configuration", func() error {
			utils.Println(added.ColoredYaml())
			return nil
		})

	form.SetOnSaveFn(func() error {
		if err := added.Validate(); err != nil {
			return err
		}
		err := added.handler.Add(added)
		if err != nil {
			return err
		}
		result = added
		return nil
	})
	form.SetOnCancelFn(func() error {
		result = added
		return nil
	})
	form.SetOnErrorFn(survey.FormShowErrorFn)
	if err := form.Render(); err != nil {
		return nil, err
	}
	return result, nil
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

func DuplicateConnector(origin *Connector, sourceManifest, targetsManifests []byte, handler ConnectorsHandler) (*Connector, error) {
	duplicated := origin.Clone(handler).
		SetTargetsManifest(targetsManifests).
		SetSourcesManifest(sourceManifest)
	if err := duplicated.askName(duplicated.Name); err != nil {
		return nil, err
	}
	if err := duplicated.askNamespace(); err != nil {
		return nil, err
	}

	checkEdit := false
	if err := survey.NewBool().
		SetKind("bool").
		SetMessage("Would you like to edit the duplicated connector before saving").
		SetRequired(true).
		SetDefault("false").
		Render(&checkEdit); err != nil {
		return nil, err
	}
	if checkEdit {
		return EditConnector(duplicated, sourceManifest, targetsManifests, handler)
	}
	if err := duplicated.Validate(); err != nil {
		return nil, err
	}
	err := duplicated.handler.Add(duplicated)
	if err != nil {
		return nil, err
	}
	return duplicated, nil
}
