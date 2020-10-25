package manager

import (
	"fmt"
	"github.com/kubemq-hub/builder/connector"
	"github.com/kubemq-hub/builder/pkg/utils"
	"github.com/kubemq-hub/builder/survey"
	"sort"
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
	sort.Slice(cm.connectors, func(i, j int) bool {
		return cm.connectors[i].Key() < cm.connectors[j].Key()
	})
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
	menu := survey.NewMenu("Select Connector to edit:").
		SetPageSize(10).
		SetDisableLoop(true).
		SetBackOption(true).
		SetErrorHandler(survey.MenuShowErrorFn)
	for _, con := range cm.connectors {
		editedCon := con.Clone(cm.handler)
		menu.AddItem(fmt.Sprintf("%s (%s)", editedCon.Key(), editedCon.Type), func() error {
			if _, err := connector.EditConnector(editedCon,
				cm.catalog.SourcesManifest,
				cm.catalog.TargetsManifest,
				cm.handler, false); err != nil {
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
func (cm *ConnectorsManager) copyConnector() error {
	cm.updateConnectors()
	menu := survey.NewMenu("Select Connector to copy:").
		SetPageSize(10).
		SetDisableLoop(true).
		SetBackOption(true).
		SetErrorHandler(survey.MenuShowErrorFn)
	for _, con := range cm.connectors {
		copiedCon := con.Clone(cm.handler)
		menu.AddItem(fmt.Sprintf("%s (%s)", copiedCon.Key(), copiedCon.Type), func() error {
			if _, err := connector.CopyConnector(copiedCon,
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
	menu := survey.NewMenu("Select Connector to delete:").
		SetPageSize(10).
		SetDisableLoop(true).
		SetBackOption(true).
		SetErrorHandler(survey.MenuShowErrorFn)
	for _, con := range cm.connectors {
		deletedCon := con.Clone(cm.handler)
		deleteFn := func() error {
			conName := deletedCon.Key()
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
				err := cm.handler.Delete(deletedCon)
				if err != nil {
					return err
				}
				utils.Println(promptConnectorDelete, conName)
			}
			return nil
		}
		menu.AddItem(fmt.Sprintf("%s (%s)", deletedCon.Key(), deletedCon.Type), deleteFn)
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
		menu.AddItem(fmt.Sprintf("%s (%s)", con.Key(), con.Type), showFn)
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
	if err := survey.NewMenu(fmt.Sprintf("Select Connectors Manager Option (Context: %s):", cm.handler.Name())).
		AddItem("<a> Add Connector", cm.addConnector).
		AddItem("<e> Edit Connector", cm.editConnector).
		AddItem("<c> Copy Connector", cm.copyConnector).
		AddItem("<d> Delete Connector", cm.deleteConnector).
		AddItem("<l> List of Connectors", cm.listConnectors).
		AddItem("<m> Catalog Management", cm.connectorsManagement).
		SetBackOption(true).
		Render(); err != nil {
		return err
	}
	return nil
}
