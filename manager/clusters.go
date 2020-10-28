package manager

import (
	"fmt"
	"github.com/kubemq-hub/builder/cluster"
	"github.com/kubemq-hub/builder/connector/common"
	"github.com/kubemq-hub/builder/pkg/uitable"
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
func (cm *ClustersManager) updateClusters() error {
	clusters, err := cm.handler.List()
	if err != nil {
		return err
	}
	var kubemqAddress []string
	for _, c := range clusters {
		kubemqAddress = append(kubemqAddress, c.EndPoints()...)
	}
	cm.loadedOptions.Add("kubemq-address", kubemqAddress)
	sort.Slice(clusters, func(i, j int) bool {
		return clusters[i].Key() < clusters[j].Key()
	})
	cm.clusters = clusters
	return nil
}
func (cm *ClustersManager) addCluster() error {
	if newCluster, err := cluster.AddCluster(
		cm.handler); err != nil {
		return fmt.Errorf("error adding new cluster: %s", err.Error())
	} else {
		utils.Println(promptClusterAdded, newCluster.Key())
	}
	return nil
}

func (cm *ClustersManager) editCluster() error {
	err := cm.updateClusters()
	if err != nil {
		return err
	}
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
				return fmt.Errorf("error editing cluster %s: %s", copiedCluster.Key(), err.Error())
			}
			utils.Println(promptClusterEdited, copiedCluster.Key())
			return nil
		})
	}
	if err := menu.Render(); err != nil {
		return err
	}
	return nil
}
func (cm *ClustersManager) copyCluster() error {
	err := cm.updateClusters()
	if err != nil {
		return err
	}
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
				return fmt.Errorf("error coping cluster %s: %s", copiedCluster.Key(), err.Error())
			}
			utils.Println(promptClusterCopied, copiedCluster.Key())
			return nil
		})
	}
	if err := menu.Render(); err != nil {
		return err
	}
	return nil
}
func (cm *ClustersManager) deleteCluster() error {
	err := cm.updateClusters()
	if err != nil {
		return err
	}
	menu := survey.NewMenu("Select Cluster to delete:").
		SetPageSize(10).
		SetDisableLoop(true).
		SetBackOption(true).
		SetErrorHandler(survey.MenuShowErrorFn)
	for _, cls := range cm.clusters {
		deletedCluster := cls.Clone(cm.handler)
		deleteFn := func() error {
			conName := deletedCluster.Key()
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
				err := cm.handler.Delete(deletedCluster)
				if err != nil {
					return fmt.Errorf("error deleting cluster %s: %s", deletedCluster.Key(), err.Error())
				}
				utils.Println(promptClusterDelete, conName)
			}
			return nil
		}
		menu.AddItem(deletedCluster.Key(), deleteFn)
	}
	if err := menu.Render(); err != nil {
		return err
	}
	return nil
}

func (cm *ClustersManager) listClusters() error {
	err := cm.updateClusters()
	if err != nil {
		return err
	}
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
	err := cm.updateClusters()
	if err != nil {
		return err
	}
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
func (cm *ClustersManager) clustersStatus() error {
	err := cm.updateClusters()
	if err != nil {
		return err
	}
	table := uitable.New()
	table.MaxColWidth = 50
	table.AddRow("NAMESPACE", "NAME", "VERSION", "STATUS", "REPLICAS", "READY", "GRPC", "REST", "API")
	for _, cls := range cm.clusters {
		table.AddRow(
			cls.Namespace,
			cls.Name,
			cls.Status.Version, cls.Status.Status, cls.Status.Replicas, cls.Status.Ready, cls.Status.Grpc, cls.Status.Rest, cls.Status.Api)
	}
	utils.Println("%s\n\n", table.String())
	return nil
}
func (cm *ClustersManager) Render() error {
	if err := survey.NewMenu(fmt.Sprintf("Select Clusters Manager Option (Context: %s):", cm.handler.Name())).
		AddItem("<a> Add Cluster", cm.addCluster).
		AddItem("<e> Edit Cluster", cm.editCluster).
		AddItem("<c> Copy Cluster", cm.copyCluster).
		AddItem("<d> Delete Cluster", cm.deleteCluster).
		AddItem("<l> List Clusters", cm.listClusters).
		AddItem("<s> Status Clusters", cm.clustersStatus).
		AddItem("<m> Cluster Integrations", cm.manageIntegrations).
		SetBackOption(true).
		SetPageSize(10).
		Render(); err != nil {
		return err
	}
	return nil
}
