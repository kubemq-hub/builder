package main

import (
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
	cm := manager.NewClustersManager(clusterHandler)
	if err := cm.Render(); err != nil {
		log.Fatal(err)
	}
}
