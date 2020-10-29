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
	loadedOptions := common.NewDefaultOptions()
	m := manager.NewManager()
	err = m.Init(loadedOptions, conHandler, clusterHandler, &Context{})
	if err != nil {
		log.Fatal(err)
	}
	if err := m.Render(); err != nil {
		log.Fatal(err)
	}

}
