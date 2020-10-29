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
	currentContext    string
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
func (cm *ClustersManager) GetClusters() ([]*cluster.Cluster, error) {
	err := cm.updateClusters()
	if err != nil {
		return nil, err
	}

	return cm.clusters, nil
}
func (cm *ClustersManager) SetCurrentContext(value string) {
	cm.currentContext = value
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
	menu := survey.NewMultiSelectMenu("Select Cluster to copy:")
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
	menu := survey.NewMultiSelectMenu("Select Cluster to delete:")
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

func (cm *ClustersManager) clustersStatus() error {
	err := cm.updateClusters()
	if err != nil {
		return err
	}
	table := uitable.New()
	table.MaxColWidth = 50
	table.AddRow("NAMESPACE", "NAME", "VERSION", "STATUS", "REPLICAS", "READY", "ADDRESS")
	if len(cm.clusters) == 0 {
		table.AddRow("<red>no clusters available</>")
	} else {
		for _, cls := range cm.clusters {
			table.AddRow(
				cls.Namespace,
				cls.Name,
				cls.Status.Version, cls.Status.Status, cls.Status.Replicas, cls.Status.Ready, cls.Status.Grpc)
		}
	}
	utils.Println("\n%s\n\n", table.String())
	return nil
}
func (cm *ClustersManager) Render() error {
	if err := survey.NewMenu(fmt.Sprintf("Select Manage Cluster Option (Context: %s):", cm.currentContext)).
		AddItem("<a> Add Cluster", cm.addCluster).
		AddItem("<e> Edit Cluster", cm.editCluster).
		AddItem("<c> Copy Clusters", cm.copyCluster).
		AddItem("<d> Delete Clusters", cm.deleteCluster).
		AddItem("<l> List Clusters", cm.listClusters).
		AddItem("<s> Status Clusters", cm.clustersStatus).
		SetBackOption(true).
		Render(); err != nil {
		return err
	}
	return nil
}
