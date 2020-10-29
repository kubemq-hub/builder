package manager

import (
	"fmt"
	"github.com/kubemq-hub/builder/cluster"
	"github.com/kubemq-hub/builder/connector"
	"github.com/kubemq-hub/builder/connector/common"
	"github.com/kubemq-hub/builder/pkg/uitable"
	"github.com/kubemq-hub/builder/pkg/utils"
	"github.com/kubemq-hub/builder/survey"
)

type IntegrationsManager struct {
	currentContext  string
	items           []*Integration
	itemsMap        map[string]*Integration
	connectorManger *ConnectorsManager
	clusterManager  *ClustersManager
}

func NewIntegrationsManager(clusterManager *ClustersManager, connectorManager *ConnectorsManager) *IntegrationsManager {
	return &IntegrationsManager{
		items:           []*Integration{},
		clusterManager:  clusterManager,
		connectorManger: connectorManager,
	}
}
func (i *IntegrationsManager) SetCurrentContext(value string) {
	i.currentContext = value
}
func (i *IntegrationsManager) populate() error {
	i.items = nil
	i.itemsMap = map[string]*Integration{}
	clustersList, err := i.clusterManager.GetClusters()
	if err != nil {
		return err
	}
	for _, cluster := range clustersList {
		clusterAddress := cluster.EndPoints()
		connectorsList, err := i.connectorManger.GetConnectors()
		if err != nil {
			return err
		}
		for _, con := range connectorsList {
			bindings := con.GetIntegrationsForCluster(clusterAddress)
			for _, binding := range bindings {
				integration := NewIntegration().
					SetBinding(binding).
					SetCluster(cluster).
					SetConnector(con)
				i.items = append(i.items, integration)
				i.itemsMap[integration.Name()] = integration
			}
		}
	}

	return nil
}

func (i *IntegrationsManager) getOrCreateConnector(cluster *cluster.Cluster, kind string) (*connector.Connector, bool, error) {
	connectors, err := i.connectorManger.GetConnectors()
	if err != nil {
		return nil, false, err
	}
	for _, c := range connectors {
		if c.Namespace == cluster.Namespace && c.Type == kind {
			return c, false, nil
		}
	}

	c := connector.NewConnector(i.connectorManger.handler, i.connectorManger.loadedOptions, i.connectorManger.catalog.TargetsManifest, i.connectorManger.catalog.SourcesManifest)
	c.Name = fmt.Sprintf("%s-%s", cluster.Name, kind)
	c.Namespace = cluster.Namespace
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
func (i *IntegrationsManager) addIntegrationWithKind(cluster *cluster.Cluster, kind string) error {
	c, isNew, err := i.getOrCreateConnector(cluster, kind)
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
		err = i.connectorManger.handler.Add(c)
	} else {
		err = i.connectorManger.handler.Edit(c)
	}
	if err != nil {
		return fmt.Errorf("error adding %s integration: %w", integration.Name, err)
	}
	utils.Println(promptIntegrationAddedConfirmation, integration.Name)
	return nil
}
func (i *IntegrationsManager) selectCluster() (*cluster.Cluster, error) {
	clustersList, err := i.clusterManager.GetClusters()
	if err != nil {
		return nil, err
	}
	var result *cluster.Cluster
	menu := survey.NewMenu("Select Cluster destination:").SetDisableLoop(true)
	for _, c := range clustersList {
		menu.AddItem(c.Key(), func() error {
			result = c
			return nil
		})
	}

	if err := menu.Render(); err != nil {
		return nil, err
	}
	return result, nil
}

func (i *IntegrationsManager) addIntegration() error {
	cluster, err := i.selectCluster()
	if err != nil {
		return err
	}
	if cluster == nil {
		return fmt.Errorf("no cluster was selected to add an integration")
	}
	menu := survey.NewMenu("Select Add Integration Type:").
		AddItem(fmt.Sprintf("Target Integration for %s", cluster.Key()), func() error {
			return i.addIntegrationWithKind(cluster, "targets")
		}).
		AddItem(fmt.Sprintf("Source Integration for %s", cluster.Key()), func() error {
			return i.addIntegrationWithKind(cluster, "sources")
		}).SetBackOption(true).
		SetErrorHandler(survey.MenuShowErrorFn).
		SetDisableLoop(true)

	if err := menu.Render(); err != nil {
		return err
	}
	return nil
}

func (i *IntegrationsManager) editIntegration() error {
	if err := i.populate(); err != nil {
		return err
	}
	menu := survey.NewMenu("Select Integration to edit:").
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
				return fmt.Errorf("error editing %s integration: %w", *integrationName, err)
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

func (i *IntegrationsManager) deleteIntegration() error {
	if err := i.populate(); err != nil {
		return err
	}
	menu := survey.NewMultiSelectMenu("Select Integration to delete:")
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
					return fmt.Errorf("error deleting %s integration: %w", cloned.Name(), err)
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

func (i *IntegrationsManager) copyIntegration() error {
	if err := i.populate(); err != nil {
		return err
	}
	menu := survey.NewMultiSelectMenu("Select Integration to copy:")
	for _, integration := range i.items {
		cloned := integration
		copyFn := func() error {
			if err := CopyIntegration(cloned, i.connectorManger); err != nil {
				return fmt.Errorf("error copying %s integration: %w", cloned.Name(), err)
			}
			utils.Println(promptIntegrationCopiedConfirmation, cloned.Name())
			return nil
		}
		menu.AddItem(cloned.Name(), copyFn)
	}
	if err := menu.Render(); err != nil {
		return err
	}
	return nil
}

func (i *IntegrationsManager) listIntegrations() error {
	if err := i.populate(); err != nil {
		return err
	}
	table := uitable.New()
	table.MaxColWidth = 50
	table.AddRow("CLUSTER", "CONNECTOR", "TYPE", "NAME", "TARGET", "SOURCE")
	if len(i.items) == 0 {
		table.AddRow("<red>no integrations available</>")
	} else {
		for _, item := range i.items {
			table.AddRow(
				item.Cluster.Key(),
				item.Connector.Key(),
				item.Connector.Type,
				item.Binding.Name,
				item.Binding.Target.Title(),
				item.Binding.Source.Title(),
			)
		}
	}
	utils.Println("\n%s\n\n", table.String())
	return nil
}

func (i *IntegrationsManager) Render() error {
	if err := survey.NewMenu(fmt.Sprintf("Select Manage Integrations Option (Context: %s):", i.currentContext)).
		AddItem("<a> Add Integration", i.addIntegration).
		AddItem("<e> Edit Integration", i.editIntegration).
		AddItem("<c> Copy Integrations", i.copyIntegration).
		AddItem("<d> Delete Integrations", i.deleteIntegration).
		AddItem("<l> List Integrations", i.listIntegrations).
		SetBackOption(true).
		SetErrorHandler(survey.MenuShowErrorFn).
		Render(); err != nil {
		return err
	}
	return nil
}
