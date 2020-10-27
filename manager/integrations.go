package manager

import (
	"fmt"
	"github.com/kubemq-hub/builder/cluster"
	"github.com/kubemq-hub/builder/connector"
	"github.com/kubemq-hub/builder/connector/common"
	"github.com/kubemq-hub/builder/pkg/utils"
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
	clusterAddress := i.cluster.EndPoints()
	for _, con := range i.connectorManger.GetConnectors() {
		bindings := con.GetIntegrationsForCluster(clusterAddress)
		for _, binding := range bindings {
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
	connectors := i.connectorManger.GetConnectors()
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
	c.Integrations = bindings
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
	takenNames := c.GetBindingNames()
	var integration *common.Binding
	switch c.Type {
	case "targets":
		integration, err = common.AddTargetIntegration(takenNames, i.connectorManger.catalog.TargetsManifest, i.connectorManger.loadedOptions)
		if err != nil {
			return err
		}
		if integration == nil {
			return nil
		}
	case "sources":
		integration, err = common.AddSourceIntegration(takenNames, i.connectorManger.catalog.SourcesManifest, i.connectorManger.loadedOptions)
		if err != nil {
			return err
		}
		if integration == nil {
			return nil
		}
	}
	if err := c.Integrations.AddIntegration(integration); err != nil {
		return err
	}
	data, err := c.Integrations.Yaml()
	if err != nil {
		return err
	}
	c.Config = string(data)

	if isNew {
		return i.connectorManger.handler.Add(c)
	}
	return i.connectorManger.handler.Edit(c)
}
func (i *Integrations) addIntegration() error {
	menu := survey.NewMenu("Select Add Integration Type:").
		AddItem(fmt.Sprintf("Target Integration for %s", i.cluster.Key()), func() error {
			return i.addIntegrationWithKind("targets")
		}).
		AddItem(fmt.Sprintf("Source Integration for %s", i.cluster.Key()), func() error {
			return i.addIntegrationWithKind("sources")
		}).SetBackOption(true).
		SetErrorHandler(survey.MenuShowErrorFn).
		SetDisableLoop(true)

	if err := menu.Render(); err != nil {
		return err
	}
	return nil
}

func (i *Integrations) editIntegration() error {
	if err := i.populate(); err != nil {
		return err
	}
	menu := survey.NewMenu("Select Integration to edit:").
		SetPageSize(15).
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
			utils.Println(promptIntegrationEditedConfirmation, *integrationName)
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
	if err := i.populate(); err != nil {
		return err
	}
	menu := survey.NewMenu("Select Integration to delete:").
		SetPageSize(15).
		SetDisableLoop(true).
		SetBackOption(true).
		SetErrorHandler(survey.MenuShowErrorFn)
	for _, integration := range i.items {
		cloned := integration
		deleteFn := func() error {
			val := false
			if err := survey.NewBool().
				SetName("confirm-delete").
				SetMessage(fmt.Sprintf("Are you sure you want to delete integration %s", cloned.Name())).
				SetRequired(true).
				SetDefault("false").
				Render(&val); err != nil {
				return err
			}
			if val {
				if err := DeleteIntegration(cloned, i.connectorManger); err != nil {
					return err
				}
				utils.Println(promptIntegrationDeleteConfirmation, cloned.Name())
			}
			return nil
		}
		menu.AddItem(cloned.Name(), deleteFn)
	}
	if err := menu.Render(); err != nil {
		return err
	}
	return nil
}

func (i *Integrations) copyIntegration() error {
	if err := i.populate(); err != nil {
		return err
	}
	menu := survey.NewMenu("Select Integration to copy:").
		SetPageSize(15).
		SetDisableLoop(true).
		SetBackOption(true).
		SetErrorHandler(survey.MenuShowErrorFn)
	for _, integration := range i.items {
		cloned := integration
		copyFn := func() error {
			if err := CopyIntegration(cloned, i.connectorManger); err != nil {
				return err
			}
			return nil
		}
		menu.AddItem(cloned.Name(), copyFn)
	}
	if err := menu.Render(); err != nil {
		return err
	}
	return nil
}

func (i *Integrations) listIntegrations() error {
	if err := i.populate(); err != nil {
		return err
	}
	menu := survey.NewMenu("Select Integration to show:").
		SetPageSize(15).
		SetDisableLoop(true).
		SetBackOption(true).
		SetErrorHandler(survey.MenuShowErrorFn)
	for _, integration := range i.items {
		cloned := integration
		showFunc := func() error {
			utils.Println("<cyan>Here is the configuration of %s Integration:</>%s", cloned.Name(), cloned.Binding.ColoredYaml())
			return nil
		}
		menu.AddItem(cloned.Name(), showFunc)
	}
	if err := menu.Render(); err != nil {
		return err
	}
	return nil
}

func (i *Integrations) Render() error {
	if err := survey.NewMenu(fmt.Sprintf("Manage Integrations for Cluster %s:", i.cluster.Key())).
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
