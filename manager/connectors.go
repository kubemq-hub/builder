package manager

import (
	"github.com/kubemq-hub/builder/connector"
	"github.com/kubemq-hub/builder/survey"
)

type ConnectorsManager struct {
	handler    ConnectorsHandler
	connectors []*connector.Connector
	catalog    *ConnectorsCatalog
}

func NewConnectorsManager(handler ConnectorsHandler) *ConnectorsManager {
	cm := &ConnectorsManager{
		handler: handler,
	}
	cm.catalog = NewConnectorCatalog()
	return cm
}

func (cm *ConnectorsManager) addConnector() error {
	if con, err := connector.NewConnector().
		SetSourcesManifest(cm.catalog.SourcesManifest).
		SetTargetsManifest(cm.catalog.TargetsManifest).
		Render(); err != nil {
		return err
	} else {
		err := cm.handler.Add(con)
		if err != nil {
			return err
		}
		return nil
	}
}
func (cm *ConnectorsManager) selectConnector() (*connector.Connector, error) {
	cm.connectors = cm.handler.List()
	selector := survey.NewListSelector("Select Connector to edit")
	for _, c := range cm.connectors {
		selector.AddItems(c)
	}
	selection, err := selector.Render()
	if err != nil {
		return nil, err
	}
	con := selection.(*connector.Connector)
	return con, nil
}
func (cm *ConnectorsManager) editConnector() error {

	return nil
}

func (cm *ConnectorsManager) deleteConnector() error {
	return nil
}

func (cm *ConnectorsManager) showConnector() error {
	return nil
}

func (cm *ConnectorsManager) listConnectors() error {
	return nil
}

func (cm *ConnectorsManager) connectorsManagement() error {
	return cm.catalog.Render()
}
func (cm *ConnectorsManager) Render() error {
	if err := cm.catalog.updateCatalog(); err != nil {
		return err
	}
	if err := survey.NewMenu("Connectors Manager: Please select").
		AddItem("Add Connector", cm.addConnector).
		AddItem("Edit Connector", cm.editConnector).
		AddItem("Delete Connector", cm.deleteConnector).
		AddItem("Show Connector", cm.showConnector).
		AddItem("List Connectors", cm.listConnectors).
		AddItem("Catalog Management", cm.connectorsManagement).
		AddItem("<-back", nil).
		Render(); err != nil {
		return err
	}
	return nil
}
