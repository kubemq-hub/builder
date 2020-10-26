package main

import (
	"github.com/kubemq-hub/builder/connector/common"
	"github.com/kubemq-hub/builder/manager"
	"log"
)

func main() {

	conHandler, err := NewConnectorsFileHandler("connectors.yaml")
	if err != nil {
		log.Fatal(err)
	}
	clusterHandler, err := NewClustersFileHandler("clusters.yaml", conHandler)
	if err != nil {
		log.Fatal(err)
	}
	catalog := manager.NewConnectorCatalog()
	_ = catalog.UpdateCatalog()
	loadedOptions := common.NewDefaultOptions()
	//	loadedOptions.Add("kubemq-address", []string{"kubemq-cluster-grpc.kubemq:50000", "kubemq-cluster-grpc.kubemq-2:50000", "Other"})
	conManager := manager.NewConnectorsManager(conHandler, catalog, loadedOptions)
	//if err := conManager.Render(); err != nil {
	//	log.Fatal(err)
	//}
	cm := manager.NewClustersManager(clusterHandler, conManager, loadedOptions)
	if err := cm.Render(); err != nil {
		log.Fatal(err)
	}
}
