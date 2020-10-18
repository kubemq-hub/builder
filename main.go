package main

import (
	"github.com/kubemq-hub/builder/connector"
	"github.com/kubemq-hub/builder/connector/common"
	"github.com/kubemq-hub/builder/pkg/utils"
	"log"
)

func main() {
	m, err := common.LoadManifestFromFile("./sources-manifest.json")
	if err != nil {
		log.Fatal(err)
	}
	//c, err := connector.NewConnector().
	//	SetSourcesManifest(m.Marshal()).
	//	Render()
	//if err != nil {
	//	log.Fatal(err)
	//}
	c, err := connector.NewConnector().
		SetSourcesManifest(m.Marshal()).
		Render()
	if err != nil {
		log.Fatal(err)
	}
	utils.Println(c.String())
}
