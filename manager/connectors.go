package manager

import (
	"fmt"
	"github.com/kubemq-hub/builder/connector"
	"github.com/kubemq-hub/builder/pkg/utils"
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
func (cm *ConnectorsManager) updateConnectors() {
	cm.connectors = cm.handler.List()
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

func (cm *ConnectorsManager) editConnector() error {
	cm.updateConnectors()
	menu := survey.NewMenu("Select Connector to edit").
		SetPageSize(10).
		SetDisableLoop(true).
		SetBackOption(true).
		SetErrorHandler(survey.MenuShowErrorFn)
	for _, con := range cm.connectors {
		editFn := func() error {
			edited := con.Clone().
				SetEditMode()
			edited, err := edited.Render()
			if err != nil {
				return err
			}
			err = cm.handler.Edit(edited)
			if err != nil {
				return err
			}
			utils.Println(promptConnectorEdit, edited.Key())
			return nil
		}
		menu.AddItem(con.Key(), editFn)
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
	menu := survey.NewMenu("Browse Connectors List").
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
		AddItem("Delete Connector", cm.deleteConnector).
		AddItem("List Connectors", cm.listConnectors).
		AddItem("Catalog Management", cm.connectorsManagement).
		SetBackOption(true).
		Render(); err != nil {
		return err
	}
	return nil
}
