package manager

import (
	"fmt"
	"github.com/kubemq-hub/builder/connector"
	"github.com/kubemq-hub/builder/pkg/utils"
	"github.com/kubemq-hub/builder/survey"
)

type ConnectorsManager struct {
	handler    connector.ConnectorsHandler
	connectors []*connector.Connector
	catalog    *ConnectorsCatalog
}

func NewConnectorsManager(handler connector.ConnectorsHandler) *ConnectorsManager {
	cm := &ConnectorsManager{
		handler: handler,
	}
	cm.catalog = NewConnectorCatalog()
	return cm
}
func (cm *ConnectorsManager) updateConnectors() {
	cm.connectors, _ = cm.handler.List()
}
func (cm *ConnectorsManager) addConnector() error {
	if _, err := connector.AddConnector(
		cm.catalog.SourcesManifest,
		cm.catalog.TargetsManifest,
		cm.handler); err != nil {
		return err
	}
	return nil
}

func (cm *ConnectorsManager) editConnector() error {
	cm.updateConnectors()
	menu := survey.NewMenu("Select Connector to edit").
		SetPageSize(10).
		SetDisableLoop(true).
		SetBackOption(true).
		SetErrorHandler(survey.MenuShowErrorFn)
	for _, con := range cm.connectors {
		menu.AddItem(con.Key(), func() error {
			if _, err := connector.EditConnector(con,
				cm.catalog.SourcesManifest,
				cm.catalog.TargetsManifest,
				cm.handler); err != nil {
				return err
			}
			return nil
		})
	}
	if err := menu.Render(); err != nil {
		return err
	}
	return nil
}
func (cm *ConnectorsManager) duplicateConnector() error {
	cm.updateConnectors()
	menu := survey.NewMenu("Select Connector to duplicate").
		SetPageSize(10).
		SetDisableLoop(true).
		SetBackOption(true).
		SetErrorHandler(survey.MenuShowErrorFn)
	for _, con := range cm.connectors {

		menu.AddItem(con.Key(), func() error {
			if _, err := connector.DuplicateConnector(con,
				cm.catalog.SourcesManifest,
				cm.catalog.TargetsManifest,
				cm.handler); err != nil {
				return err
			}
			return nil
		})
	}
	if err := menu.Render(); err != nil {
		return err
	}
	return nil
}
func (cm *ConnectorsManager) deleteConnector() error {
	cm.updateConnectors()
	menu := survey.NewMenu("Select Connector to delete").
		SetPageSize(10).
		SetDisableLoop(true).
		SetBackOption(true).
		SetErrorHandler(survey.MenuShowErrorFn)
	for _, con := range cm.connectors {
		deleteFn := func() error {
			conName := con.Key()
			val := false
			if err := survey.NewBool().
				SetName("confirm-delete").
				SetMessage(fmt.Sprintf("Are you sure you want to delete connector %s", conName)).
				SetRequired(true).
				SetDefault("false").
				Render(&val); err != nil {
				return err
			}
			if val {
				err := cm.handler.Delete(con)
				if err != nil {
					return err
				}
				utils.Println(promptConnectorDelete, conName)
			}
			return nil
		}
		menu.AddItem(con.Key(), deleteFn)
	}
	if err := menu.Render(); err != nil {
		return err
	}
	return nil
}

func (cm *ConnectorsManager) listConnectors() error {
	cm.updateConnectors()
	menu := survey.NewMenu("Browse Connectors List, Select to show configuration:").
		SetPageSize(10).
		SetBackOption(true)
	for _, con := range cm.connectors {
		str := con.ColoredYaml()
		showFn := func() error {
			utils.Println("%s\n", str)
			utils.WaitForEnter()
			return nil
		}
		menu.AddItem(con.Key(), showFn)
	}

	if err := menu.Render(); err != nil {
		return err
	}
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
		AddItem("Duplicate Connector", cm.duplicateConnector).
		AddItem("Delete Connector", cm.deleteConnector).
		AddItem("List Connectors", cm.listConnectors).
		AddItem("Catalog Management", cm.connectorsManagement).
		SetBackOption(true).
		Render(); err != nil {
		return err
	}
	return nil
}
