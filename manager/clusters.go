package manager

import (
	"fmt"
	"github.com/kubemq-hub/builder/cluster"
	"github.com/kubemq-hub/builder/connector/common"
	"github.com/kubemq-hub/builder/pkg/utils"
	"github.com/kubemq-hub/builder/survey"
	"sort"
)

type ClustersManager struct {
	handler           cluster.ClustersHandler
	clusters          []*cluster.Cluster
	connectorsManager *ConnectorsManager
	loadedOptions     common.DefaultOptions
}

func NewClustersManager(handler cluster.ClustersHandler, connectorsManager *ConnectorsManager, loadedOptions common.DefaultOptions) *ClustersManager {
	cm := &ClustersManager{
		handler:           handler,
		clusters:          nil,
		connectorsManager: connectorsManager,
		loadedOptions:     loadedOptions,
	}

	return cm
}
func (cm *ClustersManager) updateClusters() {
	cm.clusters, _ = cm.handler.List()
	var kubemqAddress []string
	for _, c := range cm.clusters {
		kubemqAddress = append(kubemqAddress, c.EndPoints()...)
	}
	cm.loadedOptions.Add("kubemq-address", kubemqAddress)
	sort.Slice(cm.clusters, func(i, j int) bool {
		return cm.clusters[i].Key() < cm.clusters[j].Key()
	})
}
func (cm *ClustersManager) addCluster() error {
	if _, err := cluster.AddCluster(
		cm.handler); err != nil {
		return err
	}
	return nil
}

func (cm *ClustersManager) editCluster() error {
	cm.updateClusters()
	menu := survey.NewMenu("Select Cluster to edit:").
		SetPageSize(10).
		SetDisableLoop(true).
		SetBackOption(true).
		SetErrorHandler(survey.MenuShowErrorFn)
	for _, c := range cm.clusters {
		copiedCluster := c.Clone(cm.handler)
		menu.AddItem(copiedCluster.Key(), func() error {
			if _, err := cluster.EditCluster(copiedCluster,
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
func (cm *ClustersManager) copyCluster() error {
	cm.updateClusters()
	menu := survey.NewMenu("Select Cluster to copy:").
		SetPageSize(10).
		SetDisableLoop(true).
		SetBackOption(true).
		SetErrorHandler(survey.MenuShowErrorFn)
	for _, c := range cm.clusters {
		copiedCluster := c.Clone(cm.handler)
		menu.AddItem(copiedCluster.Key(), func() error {
			if _, err := cluster.CopyCluster(copiedCluster,
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
func (cm *ClustersManager) deleteCluster() error {
	cm.updateClusters()
	menu := survey.NewMenu("Select Cluster to delete:").
		SetPageSize(10).
		SetDisableLoop(true).
		SetBackOption(true).
		SetErrorHandler(survey.MenuShowErrorFn)
	for _, con := range cm.clusters {
		deletedCon := con.Clone(cm.handler)
		deleteFn := func() error {
			conName := deletedCon.Key()
			val := false
			if err := survey.NewBool().
				SetName("confirm-delete").
				SetMessage(fmt.Sprintf("Are you sure you want to delete cluster %s", conName)).
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
				utils.Println(promptClusterDelete, conName)
			}
			return nil
		}
		menu.AddItem(deletedCon.Key(), deleteFn)
	}
	if err := menu.Render(); err != nil {
		return err
	}
	return nil
}

func (cm *ClustersManager) listClusters() error {
	cm.updateClusters()
	menu := survey.NewMenu("Browse Clusters List, Select to show configuration:").
		SetPageSize(10).
		SetBackOption(true)
	for _, c := range cm.clusters {
		str, _ := c.ColoredYaml()
		showFn := func() error {
			utils.Println("%s\n", str)
			utils.WaitForEnter()
			return nil
		}
		menu.AddItem(c.Key(), showFn)
	}

	if err := menu.Render(); err != nil {
		return err
	}
	return nil
}

func (cm *ClustersManager) manageIntegrations() error {
	cm.updateClusters()
	menu := survey.NewMenu("Select Cluster to manage integrations:").
		SetPageSize(10).
		SetDisableLoop(true).
		SetBackOption(true).
		SetErrorHandler(survey.MenuShowErrorFn)
	for _, c := range cm.clusters {
		clonedCluster := c.Clone(cm.handler)
		menu.AddItem(clonedCluster.Key(), func() error {
			if err := NewIntegrations(clonedCluster, cm.connectorsManager).Render(); err != nil {
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
func (cm *ClustersManager) Render() error {
	if err := survey.NewMenu(fmt.Sprintf("Select Clusters Manager Option (Context: %s):", cm.handler.Name())).
		AddItem("<a> Add Cluster", cm.addCluster).
		AddItem("<e> Edit Cluster", cm.editCluster).
		AddItem("<c> Copy Cluster", cm.copyCluster).
		AddItem("<d> Delete Cluster", cm.deleteCluster).
		AddItem("<m> Manage Cluster Integrations", cm.manageIntegrations).
		AddItem("<l> List of Clusters", cm.listClusters).
		SetBackOption(true).
		SetPageSize(10).
		Render(); err != nil {
		return err
	}
	return nil
}
