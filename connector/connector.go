package connector

import (
	"fmt"
	"github.com/kubemq-hub/builder/connector/bridges"
	"github.com/kubemq-hub/builder/connector/common"
	"github.com/kubemq-hub/builder/connector/sources"
	"github.com/kubemq-hub/builder/connector/targets"
	"github.com/kubemq-hub/builder/pkg/utils"
	"github.com/kubemq-hub/builder/survey"
	"sort"
)

type Connector struct {
	Name            string            `json:"name"`
	Namespace       string            `json:"namespace"`
	Type            string            `json:"type"`
	Replicas        int               `json:"replicas"`
	Config          string            `json:"config"`
	NodePort        int               `json:"node_port"`
	ServiceType     string            `json:"service_type"`
	Image           string            `json:"image"`
	Integrations    *common.Bindings  `json:"-" yaml:"-"`
	Bridges         *bridges.Bindings `json:"-" yaml:"-"`
	Status          *Status           `json:"-" yaml:"-"`
	loadedOptions   common.DefaultOptions
	targetManifest  []byte
	sourcesManifest []byte
	handler         ConnectorsHandler
}

func NewConnector(handler ConnectorsHandler, loadedOptions common.DefaultOptions, targetManifest, sourceManifest []byte) *Connector {
	return &Connector{
		Name:            "",
		Namespace:       "",
		Type:            "",
		Replicas:        0,
		Config:          "",
		NodePort:        0,
		ServiceType:     "",
		Image:           "",
		loadedOptions:   loadedOptions,
		targetManifest:  targetManifest,
		sourcesManifest: sourceManifest,
		handler:         handler,
		Integrations:    nil,
		Bridges:         nil,
	}
}
func (c *Connector) Clone() *Connector {
	con := &Connector{
		Name:            c.Name,
		Namespace:       c.Namespace,
		Type:            c.Type,
		Replicas:        c.Replicas,
		Config:          c.Config,
		NodePort:        c.NodePort,
		ServiceType:     c.ServiceType,
		Image:           c.Image,
		Integrations:    nil,
		Bridges:         nil,
		loadedOptions:   c.loadedOptions,
		targetManifest:  c.targetManifest,
		sourcesManifest: c.sourcesManifest,
		handler:         c.handler,
	}
	if c.Integrations != nil {
		con.Integrations = c.Integrations.Clone()
	}
	if c.Bridges != nil {
		con.Bridges = c.Bridges.Clone()
	}
	return con
}
func (c *Connector) GetBindingNames() []string {
	var list []string
	if c.Integrations != nil {
		for _, binding := range c.Integrations.Bindings {
			list = append(list, binding.Name)
		}
		return list
	}
	if c.Bridges != nil {
		for _, binding := range c.Bridges.Bindings {
			list = append(list, binding.Name)
		}
		return list
	}
	return nil
}
func (c *Connector) Update(loadedOptions common.DefaultOptions, targets, sources []byte) *Connector {
	c.loadedOptions = loadedOptions
	c.targetManifest = targets
	c.sourcesManifest = sources
	var err error
	switch c.Type {
	case "targets":
		m, _ := common.LoadManifest(c.targetManifest)
		c.Integrations, err = common.Unmarshal([]byte(c.Config))
		if err != nil {
			c.Integrations = &common.Bindings{
				Bindings: []*common.Binding{},
				Side:     "targets",
			}
		}
		c.Integrations.Side = "targets"
		c.Integrations.Update(m, c.loadedOptions)
	case "sources":
		m, _ := common.LoadManifest(c.sourcesManifest)
		c.Integrations, err = common.Unmarshal([]byte(c.Config))
		if err != nil {
			c.Integrations = &common.Bindings{
				Bindings: []*common.Binding{},
				Side:     "sources",
			}
		}
		c.Integrations.Side = "sources"
		c.Integrations.Update(m, c.loadedOptions)
	case "bridges":
		c.Bridges, err = bridges.Unmarshal([]byte(c.Config))
		if err != nil {
			c.Bridges = &bridges.Bindings{}
		}
		c.Bridges.SetDefaultOptions(c.loadedOptions)
	}

	return c
}

func (c *Connector) Key() string {
	return fmt.Sprintf("%s/%s", c.Namespace, c.Name)
}
func (c *Connector) GetManifest() *common.Manifest {
	var err error
	var m *common.Manifest
	switch c.Type {
	case "targets":
		m, err = common.LoadManifest(c.targetManifest)
	case "sources":
		m, err = common.LoadManifest(c.sourcesManifest)
	}
	if err != nil {
		return nil
	}
	return m
}

func (c *Connector) SetLoadedOptions(value common.DefaultOptions) *Connector {
	c.loadedOptions = value
	return c
}
func (c *Connector) SetHandler(value ConnectorsHandler) *Connector {
	c.handler = value
	return c
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
		return fmt.Errorf("connector replicas cannot be negative")
	}
	if c.ServiceType == "" {
		return fmt.Errorf("connector service type must have a value")
	}
	if c.NodePort < 0 {
		return fmt.Errorf("connector node-port cannot be negative")
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
		SetNamespaces(c.loadedOptions["namespaces"]).
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

func EditConnector(origin *Connector, isCopyMode bool) (*Connector, error) {
	var result *Connector
	cloned := origin.Clone()
	ftReplicas := new(string)
	*ftReplicas = fmt.Sprintf("<r> Edit Connector Replicas (%d)", cloned.Replicas)

	ftServiceType := new(string)
	*ftServiceType = fmt.Sprintf("<s> Edit Connector Service Type (%s)", cloned.ServiceType)

	ftImage := new(string)
	*ftImage = fmt.Sprintf("<d> Edit Connector Docker Image (%s)", cloned.Image)

	form := survey.NewForm(fmt.Sprintf("Select Edit %s Connector Option:", cloned.Key())).
		AddItem(fmt.Sprintf("<m> Manage Connector Bindings (%s)", cloned.Type), func() error {
			switch cloned.Type {
			case "bridges":
				bindings, err := bridges.Unmarshal([]byte(origin.Config))
				if err != nil {
					return err
				}
				cfg, err := bridges.NewBridges(cloned.Name).
					SetDefaultOptions(cloned.loadedOptions).
					SetBindings(bindings.Bindings).
					Render()
				if err != nil {
					return err
				}
				cloned.Config = string(cfg)
			case "targets":
				bindings, err := common.Unmarshal([]byte(origin.Config))
				if err != nil {
					return err
				}
				cfg, err := targets.NewTarget(cloned.Name, bindings.Bindings, cloned.loadedOptions, cloned.targetManifest).
					Render()
				if err != nil {
					return err
				}
				cloned.Config = string(cfg)
			case "sources":
				bindings, err := common.Unmarshal([]byte(origin.Config))
				if err != nil {
					return err
				}
				cfg, err := sources.NewSource(cloned.Name, bindings.Bindings, cloned.loadedOptions, cloned.sourcesManifest).
					Render()
				if err != nil {
					return err
				}
				cloned.Config = string(cfg)
			}
			utils.Println(promptConnectorContinue)
			return nil
		}).
		AddItem(ftReplicas, func() error {
			if err := cloned.askReplicas(); err != nil {
				return err
			}
			*ftReplicas = fmt.Sprintf("<r> Edit Connector Replicas (%d)", cloned.Replicas)
			return nil
		}).
		AddItem(ftServiceType, func() error {
			if err := cloned.askService(); err != nil {
				return err
			}
			*ftServiceType = fmt.Sprintf("<s> Edit Connector Service Type (%s)", cloned.ServiceType)
			return nil
		}).
		AddItem(ftImage, func() error {
			if err := cloned.askImage(); err != nil {
				return err

			}
			*ftImage = fmt.Sprintf("<d> Edit Connector Docker Image (%s)", cloned.Image)
			return nil
		}).
		AddItem("<c> Show Connector Configuration", func() error {
			utils.Println(cloned.ColoredYaml())
			return nil
		})

	form.SetOnSaveFn(func() error {
		if isCopyMode {
			result = cloned
			return nil
		}
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
	errorFn := func(err error) error {
		utils.Println("<red>error editing connector %s: %s</>", cloned.Key(), err.Error())
		return nil
	}
	form.SetOnErrorFn(errorFn)
	if err := form.Render(); err != nil {
		return nil, err
	}
	return result, nil
}

func AddConnector(handler ConnectorsHandler, loadedOptions common.DefaultOptions, targetsManifests, sourceManifest []byte) (*Connector, error) {
	var result *Connector
	added := NewConnector(handler, loadedOptions, targetsManifests, sourceManifest)
	added.Replicas = 1
	added.ServiceType = "ClusterIP"
	ftReplicas := new(string)
	*ftReplicas = fmt.Sprintf("<r> Set Connector Replicas (%d)", added.Replicas)

	ftServiceType := new(string)
	*ftServiceType = fmt.Sprintf("<s> Set Connector Service Type (%s)", added.ServiceType)

	ftImage := new(string)
	*ftImage = fmt.Sprintf("<d> Set Connector Docker Image (%s)", added.Image)

	form := survey.NewForm("Select Add Connector Option:").
		AddItem("<t> Set Connector Type and Bindings", func() error {
			utils.Println(promptConnectorStart)
			menu := survey.NewMenu("Select Connector type to add:").
				AddItem("<b> Bridges Connector", added.addBridges).
				AddItem("<t> Targets Connector", added.addTargets).
				AddItem("<s> Sources Connector", added.addSources).
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
			*ftReplicas = fmt.Sprintf("<r> Edit Connector Replicas (%d)", added.Replicas)
			return nil
		}).
		AddItem(ftServiceType, func() error {
			if err := added.askService(); err != nil {
				return err
			}
			*ftServiceType = fmt.Sprintf("<s> Edit Connector Service Type (%s)", added.ServiceType)
			return nil
		}).
		AddItem(ftImage, func() error {
			if err := added.askImage(); err != nil {
				return err

			}
			*ftImage = fmt.Sprintf("<d> Edit Connector Docker Image (%s)", added.Image)
			return nil
		}).
		AddItem("<c> Show Connector Configuration", func() error {
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
		SetDefaultOptions(c.loadedOptions).
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
	cfg, err := targets.NewTarget(c.Name, nil, c.loadedOptions, c.targetManifest).
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
	cfg, err := sources.NewSource(c.Name, nil, c.loadedOptions, c.sourcesManifest).
		Render()
	if err != nil {
		return err
	}
	c.Config = string(cfg)
	c.Type = "sources"
	return nil
}

func CopyConnector(origin *Connector) (*Connector, error) {
	copied := origin.Clone()
	if err := copied.askName(copied.Name); err != nil {
		return nil, err
	}
	if err := copied.askNamespace(); err != nil {
		return nil, err
	}

	if copied.Name == origin.Name && copied.Namespace == origin.Namespace {
		return nil, fmt.Errorf("copied connector must have different name or namespace\n")
	}
	checkEdit := false
	if err := survey.NewBool().
		SetKind("bool").
		SetMessage("Would you like to edit the copied connector before saving").
		SetRequired(true).
		SetDefault("false").
		Render(&checkEdit); err != nil {
		return nil, err
	}
	if checkEdit {
		var err error
		copied, err = EditConnector(copied, true)
		if err != nil {
			return nil, err
		}
	}
	if err := copied.Validate(); err != nil {
		return nil, err
	}
	err := copied.handler.Add(copied)
	if err != nil {
		return nil, err
	}
	return copied, nil
}

func (c *Connector) GetIntegrationsForCluster(endPoints []string) []*common.Binding {
	if c.Integrations == nil {
		return nil
	}
	var list []*common.Binding
	for _, address := range endPoints {
		list = append(list, c.Integrations.GetBindingsForCluster(address)...)
	}
	sort.Slice(list, func(i, j int) bool {
		return list[i].Name < list[j].Name
	})
	return list
}
