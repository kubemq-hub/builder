package manager

import (
	"fmt"
	"github.com/kubemq-hub/builder/cluster"
	"github.com/kubemq-hub/builder/connector"
	"github.com/kubemq-hub/builder/connector/common"
	"github.com/kubemq-hub/builder/connector/sources"
	"github.com/kubemq-hub/builder/connector/targets"
	"github.com/kubemq-hub/builder/survey"
)

type Integrations struct {
	items           []*Integration
	itemsMap        map[string]*Integration
	cluster         *cluster.Cluster
	connectorManger *ConnectorsManager
}

func NewIntegrations(cluster *cluster.Cluster, connectorManager *ConnectorsManager) *Integrations {
	return &Integrations{
		items:           []*Integration{},
		cluster:         cluster,
		connectorManger: connectorManager,
	}
}

func (i *Integrations) populate() error {
	i.items = nil
	i.itemsMap = map[string]*Integration{}
	clusterAddress := fmt.Sprintf(fmt.Sprintf("%s-grpc.%s", i.cluster.Name, i.cluster.Namespace))
	for _, con := range i.connectorManger.GetConnectors() {
		bindings := con.GetBindingsForCluster(clusterAddress)
		for _, binding := range bindings {
			binding.Print()
			integration := NewIntegration().
				SetBinding(binding).
				SetCluster(i.cluster).
				SetConnector(con)
			i.items = append(i.items, integration)
			i.itemsMap[integration.Name()] = integration
		}
	}
	return nil
}

func (i *Integrations) getOrCreateConnector(kind string) (*connector.Connector, bool, error) {
	connectors, err := i.connectorManger.handler.List()
	if err != nil {
		return nil, false, err
	}
	for _, c := range connectors {
		if c.Namespace == i.cluster.Namespace && c.Type == kind {
			return c, false, nil
		}
	}

	c := connector.NewConnector(i.connectorManger.handler, i.connectorManger.loadedOptions, i.connectorManger.catalog.TargetsManifest, i.connectorManger.catalog.SourcesManifest)
	c.Name = fmt.Sprintf("%s-%s", i.cluster.Name, kind)
	c.Namespace = i.cluster.Namespace
	c.Type = kind
	c.Replicas = 1
	c.ServiceType = "ClusterIP"

	bindings := common.NewBindings(c.Name, nil, kind, i.connectorManger.loadedOptions, c.GetManifest())
	data, err := bindings.Yaml()
	if err != nil {
		return nil, false, err
	}
	c.Config = string(data)
	return c, true, nil
}
func (i *Integrations) addIntegrationWithKind(kind string) error {
	c, isNew, err := i.getOrCreateConnector(kind)
	if err != nil {
		return err
	}
	switch c.Type {
	case "targets":
		bindings, err := common.Unmarshal([]byte(c.Config))
		if err != nil {
			return err
		}
		cfg, err := targets.NewTarget(c.Name, bindings.Bindings, i.connectorManger.loadedOptions, i.connectorManger.catalog.TargetsManifest).
			Render()
		if err != nil {
			return err
		}
		c.Config = string(cfg)
	case "sources":
		bindings, err := common.Unmarshal([]byte(c.Config))
		if err != nil {
			return err
		}
		cfg, err := sources.NewSource(c.Name, bindings.Bindings, i.connectorManger.loadedOptions, i.connectorManger.catalog.SourcesManifest).
			Render()
		if err != nil {
			return err
		}
		c.Config = string(cfg)
	}
	if isNew {
		return i.connectorManger.handler.Add(c)
	}
	return i.connectorManger.handler.Edit(c)
}
func (i *Integrations) addIntegration() error {
	form := survey.NewForm("Select Add Integration Type:").
		AddItem("Target Integration", func() error {
			return i.addIntegrationWithKind("targets")
		}).
		AddItem("Source Integration", func() error {
			return i.addIntegrationWithKind("sources")
		})

	form.SetOnSaveFn(func() error {
		return nil
	})
	form.SetOnCancelFn(func() error {
		return nil
	})
	form.SetOnErrorFn(survey.FormShowErrorFn)
	if err := form.Render(); err != nil {
		return err
	}
	return nil
}

func (i *Integrations) editIntegration() error {
	if err := i.populate(); err != nil {
		return err
	}
	menu := survey.NewMenu("Select Integration to edit:").
		SetPageSize(10).
		SetDisableLoop(true).
		SetBackOption(true).
		SetErrorHandler(survey.MenuShowErrorFn)
	for _, integration := range i.items {
		integrationName := new(string)
		*integrationName = integration.Name()
		cloned := integration
		editFunc := func() error {
			edited, err := EditIntegration(cloned, i.connectorManger)
			if err != nil {
				return err
			}
			*integrationName = edited.Name()
			return nil
		}
		menu.AddItem(*integrationName, editFunc)
	}
	if err := menu.Render(); err != nil {
		return err
	}
	return nil

}

func (i *Integrations) deleteIntegration() error {
	return nil
}

func (i *Integrations) copyIntegration() error {
	return nil
}

func (i *Integrations) listIntegrations() error {
	return nil
}

func (i *Integrations) Render() error {
	if err := survey.NewMenu(fmt.Sprintf("Manage Integrations For Cluster %s:", i.cluster.Key())).
		AddItem("<a> Add Integration", i.addIntegration).
		AddItem("<e> Edit Integration", i.editIntegration).
		AddItem("<c> Copy Integration", i.copyIntegration).
		AddItem("<d> Delete Integration", i.deleteIntegration).
		AddItem("<l> List of Integrations", i.listIntegrations).
		SetBackOption(true).
		SetErrorHandler(survey.MenuShowErrorFn).
		Render(); err != nil {
		return err
	}
	return nil
}
