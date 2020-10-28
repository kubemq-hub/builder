package manager

import (
	"fmt"
	"github.com/kubemq-hub/builder/connector"
	"github.com/kubemq-hub/builder/connector/common"
	"github.com/kubemq-hub/builder/pkg/uitable"
	"github.com/kubemq-hub/builder/pkg/utils"
	"github.com/kubemq-hub/builder/survey"
	"sort"
)

type ConnectorsManager struct {
	handler       connector.ConnectorsHandler
	connectors    []*connector.Connector
	catalog       *CatalogManager
	loadedOptions common.DefaultOptions
}

func NewConnectorsManager(handler connector.ConnectorsHandler, catalog *CatalogManager, loadedOptions common.DefaultOptions) *ConnectorsManager {
	cm := &ConnectorsManager{
		handler:       handler,
		catalog:       catalog,
		loadedOptions: loadedOptions,
	}
	return cm
}
func (cm *ConnectorsManager) SetLoadedOptions(value common.DefaultOptions) *ConnectorsManager {
	cm.loadedOptions = value
	return cm
}
func (cm *ConnectorsManager) GetConnectors() ([]*connector.Connector, error) {
	err := cm.updateConnectors()
	if err != nil {
		return nil, err
	}
	return cm.connectors, nil
}
func (cm *ConnectorsManager) updateConnectors() error {
	connectors, err := cm.handler.List()
	if err != nil {
		return err
	}
	for _, c := range connectors {
		c.Update(cm.loadedOptions, cm.catalog.TargetsManifest, cm.catalog.SourcesManifest).
			SetHandler(cm.handler)
	}
	sort.Slice(connectors, func(i, j int) bool {
		return connectors[i].Key() < connectors[j].Key()
	})
	cm.connectors = connectors
	return nil
}
func (cm *ConnectorsManager) addConnector() error {
	if newConnector, err := connector.AddConnector(
		cm.handler,
		cm.loadedOptions,
		cm.catalog.TargetsManifest,
		cm.catalog.SourcesManifest,
	); err != nil {
		return fmt.Errorf("error adding new connector: %s", err.Error())
	} else {
		utils.Println(promptConnectorAdded, newConnector.Key())
	}
	return nil
}

func (cm *ConnectorsManager) editConnector() error {
	err := cm.updateConnectors()
	if err != nil {
		return err
	}
	menu := survey.NewMenu("Select Connector to edit:").
		SetPageSize(10).
		SetDisableLoop(true).
		SetBackOption(true).
		SetErrorHandler(survey.MenuShowErrorFn)
	for _, con := range cm.connectors {
		editedCon := con.Clone()
		menu.AddItem(fmt.Sprintf("%s (%s)", editedCon.Key(), editedCon.Type), func() error {
			if _, err := connector.EditConnector(editedCon, false); err != nil {
				return fmt.Errorf("error editing connector %s: %s", editedCon.Key(), err.Error())
			}
			utils.Println(promptConnectorEdited, editedCon.Key())
			return nil
		})
	}
	if err := menu.Render(); err != nil {
		return err
	}
	return nil
}
func (cm *ConnectorsManager) copyConnector() error {
	err := cm.updateConnectors()
	if err != nil {
		return err
	}
	menu := survey.NewMenu("Select Connector to copy:").
		SetPageSize(10).
		SetDisableLoop(true).
		SetBackOption(true).
		SetErrorHandler(survey.MenuShowErrorFn)
	for _, con := range cm.connectors {
		copiedCon := con.Clone()
		menu.AddItem(fmt.Sprintf("%s (%s)", copiedCon.Key(), copiedCon.Type), func() error {
			if _, err := connector.CopyConnector(copiedCon); err != nil {
				return fmt.Errorf("error coping connector %s: %s", copiedCon.Key(), err.Error())
			}
			utils.Println(promptConnectorCopied, copiedCon.Key())
			return nil
		})
	}
	if err := menu.Render(); err != nil {
		return err
	}
	return nil
}
func (cm *ConnectorsManager) deleteConnector() error {
	err := cm.updateConnectors()
	if err != nil {
		return err
	}
	menu := survey.NewMenu("Select Connector to delete:").
		SetPageSize(10).
		SetDisableLoop(true).
		SetBackOption(true).
		SetErrorHandler(survey.MenuShowErrorFn)
	for _, con := range cm.connectors {
		deletedCon := con.Clone()
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
					return fmt.Errorf("error deleting connector %s: %s", deletedCon.Key(), err.Error())
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
	err := cm.updateConnectors()
	if err != nil {
		return err
	}
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
func (cm *ConnectorsManager) connectorsStatus() error {
	err := cm.updateConnectors()
	if err != nil {
		return err
	}
	table := uitable.New()
	table.MaxColWidth = 50
	table.AddRow("NAMESPACE", "NAME", "TYPE", "IMAGE", "REPLICAS", "STATUS")
	for _, con := range cm.connectors {
		table.AddRow(
			con.Namespace,
			con.Name,
			con.Status.Type,
			con.Status.Image,
			con.Status.Replicas,
			con.Status.Status)
	}
	utils.Println("%s\n\n", table.String())
	return nil
}
func (cm *ConnectorsManager) Render() error {
	if err := survey.NewMenu(fmt.Sprintf("Select Connectors Manager Option (Context: %s):", cm.handler.Name())).
		AddItem("<a> Add Connector", cm.addConnector).
		AddItem("<e> Edit Connector", cm.editConnector).
		AddItem("<c> Copy Connector", cm.copyConnector).
		AddItem("<d> Delete Connector", cm.deleteConnector).
		AddItem("<l> List Connectors", cm.listConnectors).
		AddItem("<s> Status Connectors", cm.connectorsStatus).
		SetBackOption(true).
		SetErrorHandler(survey.MenuShowErrorFn).
		Render(); err != nil {
		return err
	}
	return nil
}
